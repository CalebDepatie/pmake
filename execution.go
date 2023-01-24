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
	recipe_prog := fmt.Sprintf("[%v/%v]\n", recipe_num, total_recipes)
	fmt.Println(recipe_prog, "\t"+command+"\n", "\t"+out)
}

func (r *Recipe) Update(recipe_num *int) {
  // gate for if this has been chosen
  if r.Executing == nil {
    r.Executing = make(chan int)
  } else {
    _ = <- r.Executing
    return
  }

	for _, command := range r.ShellCommands {
		commandParts := strings.Split(command, " ")
		cmd := exec.Command(commandParts[0], commandParts[1:]...)

		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		RecipeFormat(command, string(stdout), *recipe_num)
	}

  *recipe_num += 1
  close(r.Executing)
}

func ExecuteGraph(cur_node *Node, recipe_num *int) {
	for _, child_node := range cur_node.Children {
		ExecuteGraph(child_node, recipe_num)
	}

	cur_node.Exec.Update(recipe_num)
}
