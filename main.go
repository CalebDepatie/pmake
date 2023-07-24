package main

import (
	"flag"
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

	// Setting up flags
	jobs := flag.Int("j", 4, "maximum number of jobs to run simultaneously")

	// Check recipe to execute
	flag.Parse()
	args := flag.Args()
	var recipe string

	if len(args) > 0 {
		recipe = args[0]

	} else {
		recipe = firstRecipe
	}

	_, ok := Project[recipe]
	if !ok {
		fmt.Println("Recipe " + recipe + " does not exist")
		return
	}

	// Create execution tree
	graphHead := CreateNode(recipe, Project)

	// Recipe execution
	x := 0
	pool := NewGoPool(*jobs)
	if ExecuteGraph(graphHead, &x, pool, nil) {
		fmt.Println("Files up to date, no work")
	}

	if failed_flag {
		panic("Error executing makefile")
	}
}
