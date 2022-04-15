package main

import (
	"flag"
	"fmt"
	"github.com/lixvbnet/sshw/sshlib"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	Name    string
	Version string
	GitHash string
)

var (
	V     = flag.Bool("v", false, "show version")
	H     = flag.Bool("h", false, "show help and exit")
	P     = flag.Int("p", 0, "port")
	U     = flag.String("u", "", "user")
	PASS  = flag.Bool("pass", false, "enter password promptly")
	T     = flag.Bool("t", false, "request terminal")
	D     = flag.Bool("d", false, "for empty values, only cover with defaults, logins section will not be used")

	templates = &promptui.SelectTemplates{
		Label:    "✨ {{ . | green}}",
		Active:   "➤ {{ .Name | cyan  }}{{if .Alias}}({{.Alias | yellow}}){{end}} {{if .Host}}{{if .User}}{{.User | faint}}{{`@` | faint}}{{end}}{{.Host | faint}}{{end}}",
		Inactive: "  {{.Name | faint}}{{if .Alias}}({{.Alias | faint}}){{end}} {{if .Host}}{{if .User}}{{.User | faint}}{{`@` | faint}}{{end}}{{.Host | faint}}{{end}}",
	}
)

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [target] [command]\n", filepath.Base(os.Args[0]))
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "options\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if !flag.Parsed() {
		flag.Usage()
		return
	}

	if *H {
		flag.Usage()
		return
	}

	if *V {
		fmt.Printf("%s version %s (%s)\n", Name, Version, GitHash)
		return
	}

	// get overrider from options(flags) [overrider has top priority]
	overrider := new(sshlib.Node)
	overrider.Port = *P
	overrider.User = *U
	if *PASS {
		prompt := promptui.Prompt{
			Label: "Enter password",
			Mask: '*',
		}
		password, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}
		overrider.Password = password
	}

	var target string
	var command string

	if len(flag.Args()) > 0 {
		target = flag.Arg(0)

		if len(flag.Args()) > 1 {
			command = flag.Arg(1)
		}
	}

	config, err := sshlib.LoadConfig(".sshw.yml")
	if err != nil {
		log.Fatal(err)
	}

	// Get node from target, or choose from menu
	var chosen *sshlib.Node
	if target == "" {
		chosen = choose(config.Nodes)
	}
	node := sshlib.GetNode(config, target, chosen, *D, overrider)
	if node == nil {
		return
	}

	// Login using node
	client := &sshlib.SSHClient{
		Host:    node.Host,
		Port:    node.Port,
	}

	err = client.LoginAuth(node.User, node.AuthMethods())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connect server ssh -p %d %s@%s version: %s\n\n", node.Port, node.User, node.Host, string(client.CurrentClient().ServerVersion()))

	// request terminal if needed
	if command == "" || *T {
		fd, state, err := client.RequestTerminal()
		defer client.RestoreTerminal(fd, state)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Shell / Start command / Run command
	if command == "" {
		if err = client.Shell(); err != nil {
			log.Fatal(err)
		}
		client.Wait()

	} else if *T {
		if err = client.Start(command); err != nil {
			log.Fatal(err)
		}
		client.Wait()

	} else {
		if _, err = client.Run(command); err != nil {
			log.Fatal(err)
		}
	}
}


func choose(nodes []*sshlib.Node) *sshlib.Node {
	prompt := promptui.Select{
		Stdout:       &BellSkipper{},
		Label:        "select host",
		Items:        nodes,
		Templates:    templates,
		Size:         20,
		HideSelected: true,
		Searcher: func(input string, index int) bool {
			node := nodes[index]
			content := fmt.Sprintf("%s %s %s", node.Name, node.User, node.Host)
			if strings.Contains(input, " ") {
				for _, key := range strings.Split(input, " ") {
					key = strings.TrimSpace(key)
					if key != "" {
						if !strings.Contains(content, key) {
							return false
						}
					}
				}
				return true
			}
			if strings.Contains(content, input) {
				return true
			}
			return false
		},
	}
	index, _, err := prompt.Run()
	if err != nil {
		return nil
	}
	return nodes[index]
}


// BellSkipper implements an io.WriteCloser that skips the terminal bell
// character (ASCII code 7), and writes the rest to os.Stderr. It is used to
// replace readline.Stdout, that is the package used by promptui to display the
// prompts.
//
// This is a workaround for the bell issue documented in
// https://github.com/manifoldco/promptui/issues/49.
type BellSkipper struct{}

// Write implements an io.WriteCloser over os.Stderr, but it skips the terminal
// bell character.
func (bs *BellSkipper) Write(b []byte) (int, error) {
	const charBell = 7 // c.f. readline.CharBell
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

// Close implements an io.WriteCloser over os.Stderr.
func (bs *BellSkipper) Close() error {
	return os.Stderr.Close()
}
