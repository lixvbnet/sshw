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

	// override fields with overrider, to make overrider.User visible to following tasks
	if overrider != nil {
		CoverDefaults(node, overrider, true)
	}

	if config != nil {
		if !defaultsOnly && config.Settings != nil {
			node.SetLogin(config.Settings.Logins, target == "") // override if target is empty (i.e. using chosen)
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

	// try as name or alias
	if !strings.Contains(target, "@") {
		node = findByNameOrAlias(config.Nodes, target)
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
			aliasNode := findByNameOrAlias(config.Nodes, node.Host)
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

func findByNameOrAlias(nodes []*Node, target string) *Node {
	target = strings.TrimSpace(target)
	for _, node := range nodes {
		if node.Name == target {
			return node
		}
		for _, alias := range strings.Split(node.Alias, ",") {
			if alias == target {
				return node
			}
		}
	}
	return nil
}
