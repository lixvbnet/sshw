package sshlib

import (
	"testing"
)

var (
	config = &Config{
		Settings: &Settings{
			Domain: "example.com",
			Logins: []*Node{
				{
					User:     "admin",
					Password: "passwordAdmin",
				},
			},
		},

		Defaults: &Node{
			User:     "default_user",
			Password: "default_password",
		},

		Nodes: []*Node{
			{
				Name:       "Node A",
				Alias:      "nodeA",
				Host:       "hostA",
				Password:   "passwordA",
			},
			{
				Name:     "Node B",
				Alias:    "nodeB",
				Host:     "hostB.com",
				User: 	  "admin",
				Password: "passwordB",
			},
			{
				Name:     "Node C",
				Alias:    "nodeC",
				Host:     "hostC",
				User: 	  "root",
				Password: "passwordC",
			},
		},
	}
)

func TestGetNode(t *testing.T) {
	var tests = []struct{
		desc string
		// input
		target string
		chosen *Node
		overrider *Node
		// want
		want Node
	}{
		// from target, without overrider
		{
			desc: "from target, without overrider",
			target: "user@host",
			want: Node{
				User: "user",
				Password: config.Defaults.Password,
				Host: "host."+config.Settings.Domain,
				Port: 22,
			},
		},
		{
			desc: "from target, without overrider",
			target: "user:pass@host",
			want: Node{
				User: "user",
				Password: "pass",
				Host: "host."+config.Settings.Domain,
				Port: 22,
			},
		},
		// from target, with overrider
		{
			desc: "from target, with overrider",
			target: "user:password@host",
			overrider: &Node{
				Password: "NewPassword",
				Port: 33,
			},
			want: Node {
				User: "user",
				Password: "NewPassword",
				Host: "host."+config.Settings.Domain,
				Port: 33,
			},
		},
		// from chosen, without overrider
		{
			desc: "from chosen, without overrider",
			chosen: config.Nodes[0],
			want: Node{
				Name: "Node A",
				Alias: "nodeA",
				User: config.Defaults.User,
				Password: "passwordA",
				Host: "hostA."+config.Settings.Domain,
				Port: 22,
			},
		},
		{
			desc: "from chosen, without overrider",
			chosen: config.Nodes[1],
			want: Node{
				Name: "Node B",
				Alias: "nodeB",
				User: "admin",
				Password: "passwordAdmin",
				Host: "hostB.com",
				Port: 22,
			},
		},
		// from chosen, with overrider
		{
			desc: "from chosen, with overrider",
			chosen: config.Nodes[0],
			overrider: &Node{
				Password: "NewPassword",
			},
			want: Node{
				Name: "Node A",
				Alias: "nodeA",
				User: config.Defaults.User,
				Password: "NewPassword",
				Host: "hostA."+config.Settings.Domain,
				Port: 22,
			},
		},
		{
			desc: "from chosen, with overrider",
			chosen: config.Nodes[2],
			overrider: &Node{
				User: "admin",
			},
			want: Node{
				Name: "Node C",
				Alias: "nodeC",
				User: "admin",
				Password: "passwordAdmin",
				Host: "hostC."+config.Settings.Domain,
				Port: 22,
			},
		},
		{
			desc: "from chosen, with overrider",
			chosen: config.Nodes[2],
			overrider: &Node{
				User: "admin",
				Password: "NewPassword",
			},
			want: Node{
				Name: "Node C",
				Alias: "nodeC",
				User: "admin",
				Password: "NewPassword",
				Host: "hostC."+config.Settings.Domain,
				Port: 22,
			},
		},
	}

	for _, test := range tests {
		if got := GetNode(config, test.target, test.chosen, test.overrider); *got != test.want {
			t.Errorf("\n//desc: %v\ntarget: %v\nchosen: %v\noverrider: %v\nwant: %v\ngot: %v\n",
				test.desc, test.target, test.chosen, test.overrider, test.want, got)
		}
	}
}
