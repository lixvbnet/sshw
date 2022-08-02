package sshlib

import "strings"

// Get node from target, or choose from menu, and then override fields with overrider.
// chosen will be ignored if target is not empty.
// If defaultsOnly is set to true, empty values will only be covered with defaults, logins section will not be used.
// Values in overrider have top priority.
func GetNode(config *Config, target string, chosen *Node, defaultsOnly bool, overrider *Node) (node *Node) {
	if target != "" {
		node = parseTarget(config, target)
	} else {
		node = chosen
	}
	if node == nil {
		return
	}

	// special case: when username is overriden, reset credentials
	if overrider != nil && overrider.User != "" && overrider.User != node.User {
		node.User = overrider.User
		node.Password = ""
	}

	if config != nil {
		if !defaultsOnly && config.Settings != nil {
			node.SetLogin(config.Settings.Logins, false)
		}
		node.SetDefaults(config.Defaults, config.Settings)
	}

	// override fields with overrider again, to make it top priority
	if overrider != nil {
		CoverDefaults(node, overrider, true)
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
	nodeAlias = strings.TrimSpace(nodeAlias)
	for _, node := range nodes {
		if node.Alias == nodeAlias {
			return node
		}
	}
	return nil
}
