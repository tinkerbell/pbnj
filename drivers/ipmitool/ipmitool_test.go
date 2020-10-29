package ipmitool

import (
	"fmt"
	"testing"

	"github.com/tinkerbell/pbnj/interfaces/bmc"
)

func TestPassThruCommand(t *testing.T) {
	for value, req := range reqs {
		s := Shell{}
		test, err := s.ComposeBmcCommand(req)
		if err != nil {
			t.Fatal(err)
		}
		if test != value {
			t.Fatal(fmt.Sprintf("got: %s; expected: %s", test, value))
		}
	}
}

var reqs = map[string]bmc.BmcRequest{
	"bmc lan set 1 ipaddr 192.168.1.1": {
		Action:  bmc.PassThruCommand,
		Command: "lan set 1 ipaddr 192.168.1.1",
	},
	"bmc reset cold": {
		Action: bmc.ColdReset,
	},
	"bmc reset warm": {
		Action: bmc.WarmReset,
	},
}
