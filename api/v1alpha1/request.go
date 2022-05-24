package update

type Request struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Device struct {
			ID       string `yaml:"id"`
			Hostname string `yaml:"hostname,omitempty"`
			Model    string `yaml:"model,omitempty"`
			Version  struct {
				Gokrazy string `yaml:"gokrazy,omitempty"`
				Kernel  string `yaml:"kernel,omitempty"`
			} `yaml:"version,omitempty"`
		} `yaml:"device"`
		Tags []struct {
			Name  string `yaml:"name"`
			Value string `yaml:"value"`
		} `yaml:"tags,omitempty"`
	} `yaml:"spec"`
}
