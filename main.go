package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("makefile")
	defer file.Close()
	if err != nil {
		fmt.Println("No makefile detected in directory")
	}

	Project, firstRecipe := GetRecipes(file)

	// Check recipe to execute
	var recipe string
	if len(os.Args[1:]) > 0 {
		recipe = os.Args[1]
	} else {
		recipe = firstRecipe
	}

	// Create execution tree
	graphHead := CreateNode(recipe, Project)

	// Recipe execution
	x := 0
	if ExecuteGraph(graphHead, &x, nil) {
		fmt.Println("Files up to date, no work")
	}

	if failed_flag {
		panic("Error executing makefile")
	}
}
