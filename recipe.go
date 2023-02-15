package main

import (
	"bufio"
	"os"
	"strings"
)

type Recipe struct {
	Dependencies  []string
	ShellCommands []string
	Executing     chan int
	Name          string
}

type envVar struct {
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

func expandVariable(input_var string, env []envVar) string {
	output_var := input_var
	for _, envVar := range env {
		var replacement string

		if len(envVar.Val) == 0 {
			replacement = " "
		} else {
			replacement = strings.Join(envVar.Val, " ")
		}

		output_var = strings.ReplaceAll(
			output_var,
			"$("+envVar.Key+")",
			replacement,
		)
	}

	return output_var
}

// handle variable and wildcard expansions
func expandProject(Project map[string]Recipe, environment []envVar) map[string]Recipe {
	// expand environment variables
	expanded_env := make([]envVar, len(environment))
	for i, env_var := range environment {
		expanded_env[i] = envVar{
			Key: env_var.Key,
			Val: make([]string, len(env_var.Val)),
		}

		for j, val := range env_var.Val {
			expanded_env[i].Val[j] = expandVariable(val, environment)
		}
	}

	// expand shell commands
	for cur_recipe, r := range Project {
		expanded_cmds := make([]string, len(r.ShellCommands))
		for i, cmd := range r.ShellCommands {
			// expand variables
			expanded_cmds[i] = expandVariable(cmd, expanded_env)
		}

		r.ShellCommands = expanded_cmds
		Project[cur_recipe] = r
	}

	return Project
}

func GetRecipes(file *os.File) (map[string]Recipe, string) {
	Project := map[string]Recipe{}

	// file parsing
	scanner := bufio.NewScanner(file)
	isRecipe := false
	curRecipe := ""
	firstRecipe := ""
	var environment []envVar

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
					environment = append(environment, envVar{
						Key: words[0],
						Val: words[2:],
					})
				}
			}
		}
	}

	Project = expandProject(Project, environment)

	return Project, firstRecipe
}
