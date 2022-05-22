## gokrazy selfupdate protocol specification

The specification protocol YAML payloads follow the Kubernetes Objects YAML schema.

### apiVersion: `update.gokrazy.org/v1alpha1`

`apiVersion` status: **`DRAFT`**

`apiVersion` date: **`DRAFT`**

1) The device's selfupdate service will accept an `-update-endpoint` flag with an HTTP/S endpoint.
1) The device's selfupdate service connects to the update-endpoint at regular intervalls to check for available updates, sending the following YAML body:

    Device selfupdate service's HTTP Update Request Body
    ```
    apiVersion: update.gokrazy.org/v1alpha1
    kind: GokrazyUpdateRequest
    metadata:
      name: "<id>-update-request-<rand>"
    spec:
      device:
        id: "<id>"
        hostname: "<hostname>"
        model: "<model>"
        version:
          gokrazy: "<BuildTimestamp>"
          kernel: "<kernel>"
      tags:
      - name: "<name>"
        value: "<value>"
    ```
1) The remote update service listening at the specified endpoint will receive the device's selfupdate request and will put together a response.
1) The remote update service will at least check for the `.spec.device.id` and compare the device against the update "database" for available updates for said device.
1) The remote update service will return an HTTP/S response for the previous request with the following YAML body:

    Remote Update service's HTTP Update Response Body
    ```
    apiVersion: update.gokrazy.org/v1alpha1
    kind: GokrazyUpdateResponse
    metadata:
      name: "<id>-update-response-<rand>"
    spec:
      device:
        id: "<id>"
      message: "<message>"
      tags:
      - name: "<name>"
        value: "<value>"
      update:
        type: "<update type>"
        links:
        - name: "<link name>"
          url: "<link url>"
        version:
          gokrazy: "<BuildTimestamp>"
          kernel: "<kernel>"
    ```
