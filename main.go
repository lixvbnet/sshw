package main

import (
	"flag"
	"fmt"
	"github.com/lixvbnet/sshlib"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"strings"
)

var (
	H     = flag.Bool("h", false, "show help and exit")
	P     = flag.Int("p", 0, "port")
	U     = flag.String("u", "", "user")
	PASS  = flag.Bool("pass", false, "enter password promptly")

	templates = &promptui.SelectTemplates{
		Label:    "✨ {{ . | green}}",
		Active:   "➤ {{ .Name | cyan  }}{{if .Alias}}({{.Alias | yellow}}){{end}} {{if .Host}}{{if .User}}{{.User | faint}}{{`@` | faint}}{{end}}{{.Host | faint}}{{end}}",
		Inactive: "  {{.Name | faint}}{{if .Alias}}({{.Alias | faint}}){{end}} {{if .Host}}{{if .User}}{{.User | faint}}{{`@` | faint}}{{end}}{{.Host | faint}}{{end}}",
	}
)

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [target]\n", os.Args[0])
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

	config, err := sshlib.LoadConfig(".sshw.yml")
	if err != nil {
		log.Fatal(err)
	}

	node := getNodeFromArgs(config)
	if node == nil {
		node = choose(config.Nodes)
	}
	if node == nil {
		return
	}

	if *P > 0 {
		node.Port = *P
	}

	if *U != "" {
		node.User = *U
	}

	if *PASS {
		prompt := promptui.Prompt{
			Label: "Enter password",
			Mask: '*',
		}
		password, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}
		node.Password = password
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
	fd, state, err := client.Shell()
	defer client.RestoreTerminal(fd, state)
	if err != nil {
		log.Fatal(err)
	}

	client.Wait()
}

func getNodeFromArgs(config *sshlib.Config) *sshlib.Node {
	if len(flag.Args()) < 1 {
		return nil
	}
	target := flag.Arg(0)

	var node *sshlib.Node

	// try as alias
	if !strings.Contains(target, "@") {
		node = findAlias(config.Nodes, target)
	}
	if node == nil {
		// login by args
		node = new(sshlib.Node)
		if !strings.Contains(target, "@") {
			node.Host = target
		} else {
			arr := strings.Split(target, "@")
			node.Host = arr[1]
			// try as alias
			aliasNode := findAlias(config.Nodes, node.Host)
			if aliasNode != nil {
				node = aliasNode
			}
			if strings.Contains(arr[0], ":") {		// user:password
				array := strings.Split(arr[0], ":")
				node.User, node.Password = array[0], array[1]
			} else {										// user
				node.User = arr[0]
			}
		}
	}
	node.SetDefaults(config.Defaults, config.Settings)
	return node
}

func findAlias(nodes []*sshlib.Node, nodeAlias string) *sshlib.Node {
	for _, node := range nodes {
		if node.Alias == nodeAlias {
			return node
		}
	}
	return nil
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
