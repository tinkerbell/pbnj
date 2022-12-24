package runner

import (
	"os"

	"gopkg.in/yaml.v2"
)

// ConfigFile is the config file.
type ConfigFile struct {
	Server Server `yaml:"server"`
	Data   `yaml:",inline"`
}

// Server is the connection details.
type Server struct {
	URL  string `yaml:"url"`
	Port string `yaml:"port"`
}

// Data of BMC resources.
type Data struct {
	Resources []Resource `yaml:"resources"`
}

// UseCases classes with names.
type UseCases struct {
	Power  []string `yaml:"power"`
	Device []string `yaml:"device"`
}

// Resource details of a single BMC.
type Resource struct {
	IP       string   `yaml:"ip"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Vendor   string   `yaml:"vendor"`
	UseCases UseCases `yaml:"useCases"`
}

// Config for the resources file.
func (c *ConfigFile) Config(name string) error {
	config, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(config, &c)
}
