package main

import (
	"bufio"
	"os"
	"strings"
)

type Recipe struct {
	Dependencies  []string
	ShellCommands []string // todo, this will require resolution
	Executing     chan int
	Name          string
}

type EnvVar struct {
	Key string
	Val []string
}

func addShellCommands(s string, cur_recipe Recipe) Recipe {
	cur_recipe.ShellCommands = append(cur_recipe.ShellCommands, strings.TrimSpace(s))
	return cur_recipe
}

func getDependencies(s string) []string {
	hasNoDepends := (s == "\n" || s == "" || s == " ")
	var depends []string

	if !hasNoDepends {
		depends = strings.Split(strings.TrimSpace(s), " ")
	}

	return depends
}

func GetRecipes(file *os.File) (map[string]Recipe, []EnvVar, string) {
	Project := map[string]Recipe{}

	// file parsing
	scanner := bufio.NewScanner(file)
	isRecipe := false
	curRecipe := ""
	firstRecipe := ""
	var environment []EnvVar

	for scanner.Scan() {
		line := scanner.Text()
		isEmpty := line == "" || line == "\n" || line == "\r\n"

		if isEmpty {
			isRecipe = false
		}

		if isRecipe {
			if firstRecipe == "" {
				firstRecipe = curRecipe
			}

			Project[curRecipe] = addShellCommands(line, Project[curRecipe])

		} else {
			if strings.Contains(line, ":") {
				var dependString string
				curRecipe, dependString, _ = strings.Cut(line, ":")

				Project[curRecipe] = Recipe{
					Dependencies: getDependencies(dependString),
					Name:         curRecipe,
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

	return Project, environment, firstRecipe
}
