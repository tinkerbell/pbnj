// +build functional

package test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/tinkerbell/pbnj/cmd"
	"gopkg.in/yaml.v2"
)

var (
	cfgData = ConfigFile{
		Server: Server{
			URL:  defaultURL,
			Port: defaultPort,
		},
	}
	defaultPort = randomInt(40000, 50000)
	defaultURL  = "localhost"
	cfg         = flag.String("config", "resources.yaml", "resources yaml file to read")
)

// ConfigFile is the config file
type ConfigFile struct {
	Server Server `yaml:"server"`
	Data   `yaml:",inline"`
}

// Server is the connection details
type Server struct {
	URL  string `yaml:"url"`
	Port string `yaml:"port"`
}

// Data of BMC resources
type Data struct {
	Resources []Resource `yaml:"resources"`
}

// UseCases classes with names
type UseCases struct {
	Power  []string `yaml:"power"`
	Device []string `yaml:"device"`
}

// Resource details of a single BMC
type Resource struct {
	IP       string   `yaml:"ip"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Vendor   string   `yaml:"vendor"`
	UseCases UseCases `yaml:"useCases"`
}

func TestMain(m *testing.M) {
	// get the resources data
	cfgData.Config(*cfg)

	if cfgData.Server.Port == defaultPort {
		// start the local internal server
		go func() {
			serverCmd := cmd.NewRootCmd()
			serverCmd.SetArgs([]string{"server", "--port", defaultPort})
			serverCmd.Execute() // nolint
		}()
	}

	time.Sleep(1 * time.Second)
	os.Exit(m.Run())
}

// Config for the resources file
func (c *ConfigFile) Config(name string) {
	flag.Parse()
	err := c.parseConfig(name)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// parseConfig reads and validates a config
func (c *ConfigFile) parseConfig(name string) error {
	// Relative to runtime directory
	_, b, _, _ := runtime.Caller(0)
	d1 := path.Join(path.Dir(b))
	config, err := ioutil.ReadFile(path.Join(d1, name))
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(config, &c)
	if err != nil {
		return err
	}

	return c.validateConfig()
}

func (c *ConfigFile) validateConfig() error {
	return nil
}

// Returns an int >= min, < max
func randomInt(min, max int) string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}
