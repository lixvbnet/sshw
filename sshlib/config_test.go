package sshlib

import (
	"fmt"
	"testing"
)

var (
	config = &Config{
		Settings: &Settings{
			Domain: "example.com",
			Logins: []*Node{
				{
					User:     "admin",
					Password: "PassAdminNew",
				},
			},
		},

		Defaults: &Node{
			Port:     0,
			User:     "root",
			Password: "PASSWORD",
		},

		Nodes: []*Node{
			{
				Name:       "nodeA",
				Alias:      "nodeA",
				Host:       "hostA",
				KeyPath:    "key_path_A",
				Passphrase: "",
			},
			{
				Name:     "nodeB",
				Alias:    "nodeB",
				Host:     "hostB",
				Password: "passB",
			},
		},
	}
)

func TestPrintConfig(t *testing.T) {
	fmt.Println(config)
}
