package sshlib

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

type Config struct {
	Settings	*Settings 	`yaml:"settings"`
	Defaults	*Node  		`yaml:"default"`
	Nodes		[]*Node  	`yaml:"nodes"`
}

func (c *Config) String() string {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
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

func (n Node) String() string {
	b, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// SetDefaults sets sensible values for unset fields in node
func (n *Node) SetDefaults(defaults *Node, settings *Settings) {
	if defaults == nil {
		defaults = new(Node)
	}
	if defaults.Port <= 0 {
		defaults.Port = 22
	}
	if defaults.User == "" {
		defaults.User = "root"
	}
	CoverDefaults(n, defaults, false)

	if !validHost(n.Host) && settings != nil && settings.Domain != "" {
		n.Host += "." + settings.Domain
	}
}

// SetLogin sets values for unset fields in node if override is set to false.
// It overrides node values if override is set to true.
func (n *Node) SetLogin(logins []*Node, override bool) {
	if logins != nil {
		for _, login := range logins {
			if n.User == login.User {
				CoverDefaults(n, login, override)
				break
			}
		}
	}
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
	// MFA auth (Keyboard Interactive)
	mfaAuth := ssh.KeyboardInteractive(func(name, instruction string, questions []string, echos []bool) ([]string, error) {
		answers := make([]string, len(questions))
		for i, question := range questions {
			if question == "Password: " && n.Password != "" {
				answers[i] = n.Password
				continue
			}
			fmt.Print(question)
			// if echo is true, display user input
			if echos[i] {
				if _, err := fmt.Scan(&answers[i]); err != nil {
					return nil, err
				}
			} else {
				answer, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println()
				if err != nil {
					return nil, err
				}
				answers[i] = string(answer)
			}
		}
		return answers, nil
	})
	authMethods = append(authMethods, mfaAuth)
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
	// executable dir
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(exe)
	for _, filename := range filenames {
		bytes, err = ioutil.ReadFile(path.Join(exeDir, filename))
		if err == nil {
			return bytes, nil
		}
	}
	return nil, err
}
