package main

import (
	"fmt"
	"os/exec"
	"strings"
)

var total_recipes int // how to keep track with parallel edges?

func init() {
	total_recipes = 0
}

func RecipeFormat(command, out string, recipe_num int) {
	recipe_prog := fmt.Sprintf("[%v/%v] : ", recipe_num, total_recipes)
	fmt.Println(recipe_prog, command+"\n", "\t"+out)
}

func (r *Recipe) Update(recipe_num *int, env []EnvVar) {
	// gate for if this has been chosen
	if r.Executing == nil {
		r.Executing = make(chan int)
	} else {
		_ = <-r.Executing
		return
	}

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

		stdout, err := cmd.Output()

    RecipeFormat(command, string(stdout), *recipe_num)

    if err != nil {
      if exiterr, ok := err.(*exec.ExitError); ok {
        fmt.Println("\t", string(exiterr.Stderr), "\n")
      } else {
        fmt.Println("\t", err.Error(), "\n")
      }
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
