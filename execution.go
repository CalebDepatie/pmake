package main

import (
	"fmt"
	"github.com/ttacon/chalk"
	"os/exec"
	"strings"
)

var total_recipes int // how to keep track with parallel edges?

func init() {
	total_recipes = 0
}

func outputHeader(recipe_num int) {
	recipe_prog := fmt.Sprintf("[%v/%v]", recipe_num, total_recipes)
	fmt.Println(chalk.Bold.TextStyle(recipe_prog))
}

func recipeFormat(command, out string) {
	fmt.Println(chalk.Bold.TextStyle("    " + command))
	fmt.Println(out)
}

func (r *Recipe) Update(recipe_num *int, env []EnvVar) {
	// gate for if this has been chosen
	if r.Executing == nil {
		r.Executing = make(chan int)
	} else {
		_ = <-r.Executing
		return
	}

	outputHeader(*recipe_num)

	for _, command := range r.ShellCommands {
		command_to_run := command
		for _, envVar := range env {
			var replacement string

			if len(envVar.Val) == 0 {
				replacement = " "
			} else {
				replacement = strings.Join(envVar.Val, " ")
			}

			command_to_run = strings.ReplaceAll(
				command_to_run,
				"$("+envVar.Key+")",
				replacement,
			)
		}
		// fmt.Println(command_to_run)
		commandParts := strings.Split(command_to_run, " ")
		cmd := exec.Command(commandParts[0], commandParts[1:]...)

		stdout, err := cmd.CombinedOutput()

		recipeFormat(command_to_run, string(stdout))

		if err != nil {
			fmt.Println(err.Error(), "\n")
		}
	}

	*recipe_num += 1
	close(r.Executing)
}

func ExecuteGraph(cur_node *Node, recipe_num *int, env []EnvVar) {
	for _, child_node := range cur_node.Children {
		ExecuteGraph(child_node, recipe_num, env)
	}

	if !cur_node.Executed {
		cur_node.Exec.Update(recipe_num, env)
		cur_node.Executed = true
	}
}
