// +build functional

package test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/tinkerbell/pbnj/cmd"
	"gopkg.in/yaml.v2"
)

var (
	cfgData Data
	port    = "40041"
	cfg     = flag.String("config", "resources.yaml", "resources yaml file to read")
)

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
	cfgData = Config()
	// start the server
	go func() {
		serverCmd := cmd.NewRootCmd()
		serverCmd.SetArgs([]string{"server", "--port", port})
		serverCmd.Execute() // nolint
	}()
	time.Sleep(1 * time.Second)
	os.Exit(m.Run())
}

// Config for the resources file
func Config() Data {
	flag.Parse()
	cfgData, err := parseConfig(*cfg)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return cfgData
}

// parseConfig reads and validates a config
func parseConfig(name string) (Data, error) {
	c := Data{}
	// Relative to runtime directory
	_, b, _, _ := runtime.Caller(0)
	d1 := path.Join(path.Dir(b))
	config, err := ioutil.ReadFile(path.Join(d1, name))
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(config, &c)
	if err != nil {
		return c, err
	}

	return c, validateConfig(c)
}

func validateConfig(c Data) error {
	return nil
}
