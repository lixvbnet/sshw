package sshlib

import (
	"fmt"
	"testing"
)

var (
	settings = &Settings{Domain: "example.com"}

	defaults = &Node{
		Name:           "",
		Alias:          "",
		Host:           "",
		Port:           0,
		User:           "root",
		Password:       "PASSWORD",
		KeyPath:        "",
		Passphrase:     "",
	}

	node = &Node{
		Name:           "dev_local",
		Alias:          "dev",
		Host:           "hostA",
		Port:           0,
		User:           "",
		Password:       "pass",
		KeyPath:        "some_key_path",
		Passphrase:     "",
	}
)

func TestSetDefaults(t *testing.T) {
	node.SetDefaults(defaults, settings)

	fmt.Println(node)
	if node.User != "root" || node.Password != "pass" {
		t.Errorf("Error: node.User=%s, node.Password=%s", node.User, node.Password)
	}
}

func TestLoadConfig(t *testing.T) {
	var config *Config
	config, err := LoadConfig("fileNotExist.yml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config)

	config, err = LoadConfig("config_example.yml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config)
}
