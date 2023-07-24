package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Recipe struct {
	Dependencies  []string
	ShellCommands []string
	Executing     chan int
	Name          string
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

func expandVariable(input_var string, env map[string][]string) string {
	output_var := input_var
	for key, val := range env {
		var replacement string

		if len(val) == 0 {
			replacement = " "
		} else {
			replacement = strings.Join(val, " ")
		}

		output_var = strings.ReplaceAll(
			output_var,
			"$("+key+")",
			replacement,
		)
	}

	return output_var
}

func expandWildcards(s string) string {
	string_parts := strings.Split(s, " ")

	for i, cur_string := range string_parts {
		if !strings.ContainsRune(cur_string, '*') {
			continue
		}

		matches, err := filepath.Glob(cur_string)
		if err != nil {
			fmt.Println("Error expanding " + cur_string + ": " + err.Error())
		}

		string_parts[i] = strings.Join(matches, " ")
	}

	return strings.Join(string_parts, " ")
}

// handle variable and wildcard expansions
func expandProject(Project map[string]Recipe, environment map[string][]string) map[string]Recipe {
	// expand environment variables
	expanded_env := make(map[string][]string)
	for key, val := range environment {
		expanded_env[key] = make([]string, len(val))

		for i, v := range val {
			expanded_env[key][i] = expandVariable(v, environment)
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

		for i, dep := range r.Dependencies {
			r.Dependencies[i] = expandWildcards(dep)
		}
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

	environment := make(map[string][]string)
	environment[".RECIPEPREFIX"] = []string{""}
	environment[".DEFAULT_GOAL"] = []string{""}

	for scanner.Scan() {
		line := scanner.Text()
		isEmpty := line == "" || line == "\n" || line == "\r\n"

		if isEmpty {
			isRecipe = false
		}

		if isRecipe {
			// set first recipe for determining which to execute first
			if firstRecipe == "" {
				firstRecipe = curRecipe
			}

			// ensure that it begins with .RECIPEPREFIX
			if strings.HasPrefix(line, strings.Join(environment[".RECIPEPREFIX"], " ")) {
				// remove prefix
				line = strings.TrimPrefix(line, strings.Join(environment[".RECIPEPREFIX"], " "))
				Project[curRecipe] = addShellCommands(line, Project[curRecipe])
			
			} else {
				isRecipe = false
			}

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
				name, values, _ := strings.Cut(line, "=")

				name = strings.TrimSpace(name)
				values = strings.TrimSpace(values)
				words := strings.Split(values, " ")

				environment[name] = words
			}
		}
	}

	Project = expandProject(Project, environment)

	if environment[".DEFAULT_GOAL"][0] != "" {
		firstRecipe = strings.Join(environment[".DEFAULT_GOAL"], " ")
	}

	return Project, firstRecipe
}
