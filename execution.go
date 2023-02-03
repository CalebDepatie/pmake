package main

import (
	"fmt"
	"github.com/ttacon/chalk"
	"os/exec"
	"strings"
)

var total_recipes int // how to keep track with parallel edges?
var failed_flag bool  // a flag to track if any recipe failes

func init() {
	total_recipes = 0
	failed_flag = false
}

func outputHeader(recipe_num int) string {
	recipe_prog := fmt.Sprintf("[%v/%v]", recipe_num, total_recipes)
	return chalk.Bold.TextStyle(recipe_prog)
}

func recipeFormat(command, out string) string {
	s := chalk.Bold.TextStyle("    " + command)
	s += "\n" + out + "\n"
	return s
}

func createCommand(command string, env []EnvVar) string {
	for _, envVar := range env {
		var replacement string

		if len(envVar.Val) == 0 {
			replacement = " "
		} else {
			replacement = strings.Join(envVar.Val, " ")
		}

		command = strings.ReplaceAll(
			command,
			"$("+envVar.Key+")",
			replacement,
		)
	}

	return command
}

func (r *Recipe) update(recipe_num *int, env []EnvVar) {
	// gate for if this has been chosen
	if r.Executing == nil {
		r.Executing = make(chan int)
		*recipe_num += 1
	} else {
		_ = <-r.Executing
		return
	}

	output_string := outputHeader(*recipe_num) + "\n"

	// check if the work actually needs to be done

	for _, command := range r.ShellCommands {
		command_to_run := createCommand(command, env)

		commandParts := strings.Split(command_to_run, " ")
		cmd := exec.Command(commandParts[0], commandParts[1:]...)

		stdout, err := cmd.CombinedOutput()

		output_string += recipeFormat(command_to_run, string(stdout))

		if err != nil {
			output_string += err.Error()
			failed_flag = true
		}
	}

	fmt.Printf(output_string)

	close(r.Executing)
}

func ExecuteGraph(cur_node *Node, recipe_num *int, env []EnvVar) {
	for _, child_node := range cur_node.Children {
		ExecuteGraph(child_node, recipe_num, env)
	}

	if !cur_node.Executed {
		cur_node.Exec.update(recipe_num, env)
		cur_node.Executed = true
	}
}
