package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/AliBa1/gator/internal/config"
	"github.com/AliBa1/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	state := state{
		config: &c,
	}
	commands := NewCommands()

	db, err := sql.Open("postgres", state.config.DbURL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}
	dbQueries := database.New(db)
	state.database = dbQueries

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

	command := command{
		name: userCommandName,
		args: args,
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
	// c, err = config.Read()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// c.Print()
}
