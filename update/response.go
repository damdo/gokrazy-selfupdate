package update

type Response struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Device struct {
			ID string `yaml:"id"`
		} `yaml:"device"`
		Message string `yaml:"message,omitempty"`
		Tags    []struct {
			Name  string `yaml:"name"`
			Value string `yaml:"value"`
		} `yaml:"tags,omitempty"`
		Update struct {
			Type  string `yaml:"type"`
			Links []struct {
				Name string `yaml:"name"`
				URL  string `yaml:"url"`
			} `yaml:"links"`
			Version struct {
				Gokrazy string `yaml:"gokrazy"`
				Kernel  string `yaml:"kernel"`
			} `yaml:"version"`
		} `yaml:"update"`
	} `yaml:"spec"`
}
