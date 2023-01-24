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
	Executing     chan int
}

type EnvVar struct {
	Key string
	Val []string
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
	firstRecipe := ""
	var environment []EnvVar
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line == "\n" || line == "\r\n" {
			isRecipe = false
		}

		if isRecipe {
			if firstRecipe == "" {
				firstRecipe = curRecipe
			}

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
			} else if strings.Contains(line, "=") {
				words := strings.Split(line, " ")

				if words[1] == "=" {
					environment = append(environment, EnvVar{
						Key: words[0],
						Val: words[2:],
					})
				}
			}
		}
	}

	// Create execution tree
	graphHead := CreateNode(firstRecipe, Project)

	// Recipe execution
	x := 1
	ExecuteGraph(graphHead, &x, environment)
}
