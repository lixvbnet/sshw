package sshlib

import "strings"

// Get node from target, or choose from menu, and then override fields with overrider
// chosen will be ignored if target is not empty
func GetNode(config *Config, target string, chosen *Node, overrider *Node) (node *Node) {
	if target != "" {
		node = parseTarget(config, target)
	} else {
		node = chosen
	}
	if node == nil {
		return
	}

	if overrider != nil {
		CoverDefaults(node, overrider, true)
	}
	if config != nil {
		if config.Settings != nil {
			node.SetLogin(config.Settings.Logins, target == "") // override if target is empty (i.e. using chosen)
		}
		node.SetDefaults(config.Defaults, config.Settings)
	}
	return node
}

func parseTarget(config *Config, target string) *Node {
	if target == "" {
		return nil
	}

	var node *Node

	// try as alias
	if !strings.Contains(target, "@") {
		node = findAlias(config.Nodes, target)
	}
	if node == nil {
		// login by args
		node = new(Node)
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
	return node
}

func findAlias(nodes []*Node, nodeAlias string) *Node {
	for _, node := range nodes {
		if node.Alias == nodeAlias {
			return node
		}
	}
	return nil
}
