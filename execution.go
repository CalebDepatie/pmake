package main

import (
	"fmt"
	"github.com/ttacon/chalk"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	total_recipes int  // how to keep track with parallel edges?
	failed_flag   bool // a flag to track if any recipe fails
	filesystem    fs.FS
)

func init() {
	total_recipes = 0
	failed_flag = false

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Couldn't get working directory: " + err.Error())
	}

	filesystem = os.DirFS(wd)
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

func (r *Recipe) update(recipe_num *int) {
	// gate for if this has been chosen
	if r.Executing == nil {
		r.Executing = make(chan int)
		*recipe_num += 1
	} else {
		_ = <-r.Executing
		return
	}

	output_string := outputHeader(*recipe_num) + "\n"

	for _, command := range r.ShellCommands {

		commandParts := strings.Split(command, " ")
		cmd := exec.Command(commandParts[0], commandParts[1:]...)

		stdout, err := cmd.CombinedOutput()

		output_string += recipeFormat(command, string(stdout))

		if err != nil {
			output_string += err.Error()
			failed_flag = true
		}
	}

	fmt.Printf(output_string)

	close(r.Executing)
}

func checkFile(node Node) bool {
	result_file := node.Exec.Name
	src_files := node.Exec.Dependencies

	result_info, err := fs.Stat(filesystem, result_file)
	if err != nil {
		// fmt.Println("Debug: could not stat result file")
		return false
	}

	result_time := result_info.ModTime()

	for _, file := range src_files {
		file_info, err := fs.Stat(filesystem, file)
		if err != nil {
			// fmt.Println("Debug: could not stat file " + file)
			return false
		}

		file_time := file_info.ModTime()

		if file_time.After(result_time) {
			return false
		}
	}

	return true
}

func ExecuteGraph(cur_node *Node, recipe_num *int, parent_wait *sync.WaitGroup) bool {

	notifyParent := func() {
		if parent_wait != nil {
			parent_wait.Done()
		}
	}

	child_wait := new(sync.WaitGroup)
	child_wait.Add(len(cur_node.Children))
	for _, child_node := range cur_node.Children {
		go ExecuteGraph(child_node, recipe_num, child_wait)
	}

	child_wait.Wait()

	skipWork := checkFile(*cur_node)

	if skipWork {
		notifyParent()
		return skipWork
	}

	if !cur_node.Executed {
		cur_node.Exec.update(recipe_num)
		cur_node.Executed = true
	}

	notifyParent()

	return skipWork
}
