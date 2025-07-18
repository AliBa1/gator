package main

import (
	"fmt"
	"os"

	"github.com/AliBa1/gator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	state := config.State{
		Config: &c,
	}
	commands := config.NewCommands()

	if len(os.Args) < 2 {
		fmt.Printf("Enter a command to use the app")
		os.Exit(1)
		return
	}

	var userCommandName string
	userCommandName = os.Args[1]

	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	command := config.Command{
		Name: userCommandName,
		Args: args,
	}
	err = commands.Run(&state, command)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	// err = c.SetUser("ali")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	//
	c, err = config.Read()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c.Print()
}
