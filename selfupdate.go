package main

import (
	"archive/zip"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gokrazy/updater"
	log "github.com/sirupsen/logrus"

	updateapi "github.com/damdo/gokrazy-selfupdate/api/v1alpha1"
	yaml "gopkg.in/yaml.v3"
)

// TODO: compute this
var deviceID = "myID"
var buildTimestamp = "myBuildTimestamp"

var updateEndpoint, checkInterval string

func selfupdate(ctx context.Context, response *updateapi.Response) error {
	log.Println("starting self-update procedure")

	var rootReader, bootReader, mbrReader io.ReadCloser

	switch response.Spec.Update.Type {
	case "zip-multi-part":
		link := response.Spec.Update.Links[0]
		if link.Name != "drive" {
			return fmt.Errorf("unrecognized link name `%s`, for link type `zip-multi-part`, name must be `drive`", link.Name)
		}

		log.Println("downloading update file")
		filePath := "/tmp/drive.zip"
		err := downloadFile(filePath, link.URL)
		if err != nil {
			return fmt.Errorf("unable to download update file: %w", err)
		}

		log.Println("loading new disk partitions from update file")
		r, err := zip.OpenReader(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()

		for _, f := range r.File {
			switch f.Name {
			case "mbr.img":
				mbrReader, err = f.Open()
				if err != nil {
					log.Fatal(err)
				}
			case "boot.img":
				bootReader, err = f.Open()
				if err != nil {
					log.Fatal(err)
				}
			case "root.squashfs":
				rootReader, err = f.Open()
				if err != nil {
					log.Fatal(err)
				}
			}
		}

	default:
		return fmt.Errorf("unrecognized .spec.update.type")
	}

	httpPassword, err := readConfigFile("gokr-pw.txt")
	if err != nil {
		return fmt.Errorf("could read neither /perm/gokr-pw.txt, nor /etc/gokr-pw.txt, nor /gokr-pw.txt: %w", err)
	}

	httpPort, err := readConfigFile("http-port.txt")
	if err != nil {
		return fmt.Errorf("could read neither /perm/http-port.txt, nor /etc/http-port.txt, nor /http-port.txt: %w", err)
	}

	uri := fmt.Sprintf("http://gokrazy:%s@127.0.0.1:%s/", httpPassword, httpPort)

	log.Println("checking target partuuid support")
	target, err := updater.NewTarget(uri, http.DefaultClient)
	if err != nil {
		log.Fatalf("checking target partuuid support: %v", err)
	}

	// Start with the root file system because writing to the non-active
	// partition cannot break the currently running system.
	log.Println("updating root file system")
	if err := target.StreamTo("root", rootReader); err != nil {
		log.Fatalf("updating root file system: %v", err)
	}
	rootReader.Close()

	log.Println("updating boot file system")
	if err := target.StreamTo("boot", bootReader); err != nil {
		log.Fatalf("updating boot file system: %v", err)
	}
	bootReader.Close()

	// Only relevant when running on non-Raspberry Pi devices.
	// As it does not use an MBR.
	log.Println("updating MBR")
	if err := target.StreamTo("mbr", mbrReader); err != nil {
		log.Fatalf("updating MBR: %v", err)
	}
	mbrReader.Close()

	log.Println("switching to non-active partition")
	if err := target.Switch(); err != nil {
		log.Fatalf("switching to non-active partition: %v", err)
	}

	log.Println("reboot")
	if err := target.Reboot(); err != nil {
		log.Fatalf("reboot: %v", err)
	}

	return nil
}

func checkForUpdates(ctx context.Context) (*updateapi.Response, error) {
	body := updateapi.Request{
		APIVersion: "update.gokrazy.org/v1alpha1",
		Kind:       "GokrazyUpdateRequest",
	}
	body.Metadata.Name = "name"
	body.Spec.Device.ID = deviceID

	reqBody, err := yaml.Marshal(&body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", updateEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/yaml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var response updateapi.Response
	err = yaml.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func shouldUpdate(response *updateapi.Response) bool {

	if response.Spec.Device.ID != deviceID {
		log.Errorf("update response's device id: %s, differs from the actual device id: %s, aborting", response.Spec.Device.ID, deviceID)
		return false
	}

	if response.Spec.Update.Version.Gokrazy == buildTimestamp {
		log.Infof("device's gokrazy version: %s is already the desired one, skipping", response.Spec.Update.Version.Gokrazy)
		return false
	}

	log.Infof("device's gokrazy version: %s, desired version: %s, proceeding with the update", buildTimestamp, response.Spec.Update.Version.Gokrazy)

	return true
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Info("gokrazy's selfupdate service starting up..")

	flag.StringVar(&updateEndpoint, "update-endpoint", "", "the HTTP/S endpoint of the update service")
	flag.StringVar(&checkInterval, "check-interval", "1h", "the time duration interval between checks to the update service. default: 1h")
	flag.Parse()

	if updateEndpoint == "" {
		log.Fatalln("flag --update-endpoint must be provided")
	}

	interval, err := time.ParseDuration(checkInterval)
	if err != nil {
		log.Fatalln(err)
	}

	log.Infof("entering update checking loop with interval: %s", interval.String())
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			log.Info("stopping update checking")
			ticker.Stop()
			return
		case <-ticker.C:
			response, err := checkForUpdates(ctx)
			if err != nil {
				log.Errorf("unable to check for updates: %v", err)
				continue
			}

			if shouldUpdate(response) {
				if err := selfupdate(ctx, response); err != nil {
					log.Fatal(err)
				}
				// the update is now correctly written to the disk partitions
				// and the reboot is in progress
				// sleep until the context chan is closed, then exit cleanly
				<-ctx.Done()
				os.Exit(0)
			}
		}
	}
}
