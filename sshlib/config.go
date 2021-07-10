package sshlib

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
	"path"
	"strings"
)

type Config struct {
	Settings	*Settings 	`yaml:"settings"`
	Defaults	*Node  		`yaml:"default"`
	Nodes		[]*Node  	`yaml:"nodes"`
}

// SetDefaults sets sensible values for unset fields in each node
func (c *Config) SetDefaults() {
	for _, node := range c.Nodes {
		node.SetDefaults(c.Defaults, c.Settings)
	}
}

type Settings struct {
	Domain		string		`yaml:"domain"`
	Logins		[]*Node		`yaml:"logins"`
}

type Node struct {
	Name           string           `yaml:"name"`
	Alias          string           `yaml:"alias"`
	Host           string           `yaml:"host"`
	Port           int              `yaml:"port"`
	User           string           `yaml:"user"`
	Password       string           `yaml:"password"`
	KeyPath        string           `yaml:"keypath"`
	Passphrase     string           `yaml:"passphrase"`
}

// SetDefaults sets sensible values for unset fields in node
func (n *Node) SetDefaults(defaults *Node, settings *Settings) {
	if !validHost(n.Host) && settings != nil && settings.Domain != "" {
		n.Host += "." + settings.Domain
	}
	if n.Port <= 0 {
		if defaults != nil && defaults.Port > 0 {
			n.Port = defaults.Port
		} else {
			n.Port = 22
		}
	}
	if n.User == "" {
		if defaults != nil && defaults.User != "" {
			n.User = defaults.User
		} else {
			n.User = "root"
		}
	}
	if settings != nil && settings.Logins != nil {
		for _, login := range settings.Logins {
			if n.User == login.User {
				CoverDefaults(n, login)
				break
			}
		}
	}
	if n.Password == "" {
		if defaults != nil && defaults.Password != "" {
			n.Password = defaults.Password
		}
	}
	CoverDefaults(n, defaults)
}

func validHost(host string) bool {
	host = strings.TrimSpace(host)
	validHosts := []string{"localhost", "127.0.0.1"}
	for _, h := range validHosts {
		if host == h {
			return true
		}
	}
	return strings.Contains(host, ".") || strings.Contains(host, ":")
}

func (n *Node) AuthMethods() (authMethods []ssh.AuthMethod) {
	// RSA auth
	rsaAuth, err := n.rsaAuth()
	if err != nil {
		fmt.Println(err)
	}
	if rsaAuth != nil {
		authMethods = append(authMethods, rsaAuth)
	}
	// Password auth
	if n.Password != "" {
		authMethods = append(authMethods, ssh.Password(n.Password))
	}
	return authMethods
}

func (n *Node) rsaAuth() (rsaAuth ssh.AuthMethod, err error) {
	var pemBytes []byte

	if n.KeyPath == "" {
		u, err := user.Current()
		if err != nil {
			return nil, nil	// ignore the error
		}
		pemBytes, err = ioutil.ReadFile(path.Join(u.HomeDir, ".ssh/id_rsa"))
		if err != nil {
			return nil, nil	// possibly not exist, ignore
		}
	} else {
		pemBytes, err = ioutil.ReadFile(n.KeyPath)
		if err != nil {
			return nil, err
		}
	}

	var signer ssh.Signer
	if n.Passphrase == "" {
		signer, err = ssh.ParsePrivateKey(pemBytes)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(n.Passphrase))
	}
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func LoadConfig(filenames ...string) (config *Config, err error) {
	b, err := loadConfigBytes(filenames)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return
	}
	config.SetDefaults()
	return config, nil
}

func loadConfigBytes(filenames []string) (bytes []byte, err error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	// homedir
	for _, filename := range filenames {
		bytes, err = ioutil.ReadFile(path.Join(u.HomeDir, filename))
		if err == nil {
			return bytes, nil
		}
	}
	// relative
	for _, filename := range filenames {
		bytes, err = ioutil.ReadFile(filename)
		if err == nil {
			return bytes, nil
		}
	}
	return nil, err
}
