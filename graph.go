package main

type Node struct {
	Exec     Recipe
	Children []*Node
	Executed bool
}

var nodes map[string]*Node

func init() {
	nodes = make(map[string]*Node)
}

func CreateNode(recipe_name string, proj map[string]Recipe) *Node {

	node, exists := nodes[recipe_name]
	if exists {
		return node
	}

	var depends []*Node
	for _, depend_name := range proj[recipe_name].Dependencies {
		if _, exists := proj[depend_name]; exists { // todo: track these in a different way
			depends = append(depends, CreateNode(depend_name, proj))
		}
	}

	total_recipes += 1
	new_node := Node{
		Exec:     proj[recipe_name],
		Children: depends,
		Executed: false,
	}

	nodes[recipe_name] = &new_node

	return &new_node
}
