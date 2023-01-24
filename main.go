package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Recipe struct {
	Dependencies  []string
	ShellCommands []string // todo, this will require resolution
}

func main() {

	file, err := os.Open("makefile")
	defer file.Close()
	if err != nil {
		fmt.Println("No makefile detected in directory")
	}

	Project := map[string]Recipe{}

	// file parsing
	scanner := bufio.NewScanner(file)
	isRecipe := false
	curRecipe := ""
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line == "\n" || line == "\r\n" {
			isRecipe = false
		}

		if isRecipe {
			workingRecipe := Project[curRecipe]
			workingRecipe.ShellCommands = append(workingRecipe.ShellCommands, strings.TrimSpace(line))
			Project[curRecipe] = workingRecipe

		} else {
			if strings.Contains(line, ":") {
				var dependString string
				curRecipe, dependString, _ = strings.Cut(line, ":")

				hasNoDepends := (dependString == "\n" || dependString == "" || dependString == " ")
				var depends []string

				if hasNoDepends {

				} else {
					depends = strings.Split(strings.TrimSpace(dependString), " ")
				}

				Project[curRecipe] = Recipe{
					Dependencies: depends,
				}
				isRecipe = true
			}

		}
	}

	// Create execution tree
	graphHead := CreateNode("test", Project)

	// Recipe execution
	x := 0
	ExecuteGraph(graphHead, &x)
}
