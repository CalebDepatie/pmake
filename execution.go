package main

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
