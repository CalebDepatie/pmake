package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Recipe struct {
	Dependencies  []string
	ShellCommands []string // todo, this will require resolution
}

type Node struct {
	Exec     Recipe
	Children []Node
}

func CreateNode(recipe_name string, proj map[string]Recipe) Node {
  var depends []Node
  for _, depend_name := range proj[recipe_name].Dependencies {
    depends = append(depends, CreateNode(depend_name, proj))
  }

  return Node{
    Exec:     proj[recipe_name],
    Children: depends,
  }
}

func ExecuteGraph(cur_node Node) {
  for _, child_node := range cur_node.Children {
    ExecuteGraph(child_node)
  }

  for _, command := range cur_node.Exec.ShellCommands {
    commandParts := strings.Split(command, " ")
    cmd := exec.Command(commandParts[0], commandParts[1:]...)

    stdout, err := cmd.Output()

    if err != nil {
      fmt.Println(err.Error())
      return
    }

    fmt.Println(string(stdout))
  }
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
				Project[curRecipe] = Recipe{
					Dependencies: strings.Split(strings.TrimSpace(dependString), " "),
				}
				isRecipe = true
			}

		}
	}

	// Create execution tree
	graphHead := CreateNode("test", Project)
	fmt.Println(graphHead)
  fmt.Println("")

  // Recipe execution
  ExecuteGraph(graphHead)
}
