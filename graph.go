package main

type Node struct {
	Exec     Recipe
	Children []Node
}

func CreateNode(recipe_name string, proj map[string]Recipe) Node {
	var depends []Node
	for _, depend_name := range proj[recipe_name].Dependencies {
		depends = append(depends, CreateNode(depend_name, proj))
	}

	return Node{
		Exec:     proj[recipe_name],
		Children: depends,
	}
}
