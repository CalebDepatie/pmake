package main

type Node struct {
	Exec     Recipe
	Children []*Node
}

var nodes map[string]*Node

func init() {
	nodes = make(map[string]*Node)
}

func CreateNode(recipe_name string, proj map[string]Recipe) *Node {
	var depends []*Node
	for _, depend_name := range proj[recipe_name].Dependencies {
		depends = append(depends, CreateNode(depend_name, proj))
	}

	node, exists := nodes[recipe_name]
	if exists {
		return node
	}

	total_recipes += 1
	new_node := Node{
		Exec:     proj[recipe_name],
		Children: depends,
	}

	nodes[recipe_name] = &new_node

	return &new_node
}
