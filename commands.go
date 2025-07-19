package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AliBa1/gator/internal/config"
	"github.com/AliBa1/gator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	config   *config.Config
	database *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func NewCommands() commands {
	commands := commands{
		commands: make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)

	return commands
}

// func CommandsMap() map[string]func(*State, Command) error {
// 	commands := make(map[string]func(*State, Command) error)
// 	commands["login"] = handlerLogin
// 	return commands
// }

func (c *commands) Run(s *state, cmd command) error {
	handlerFunc, ok := c.commands[cmd.name]
	if !ok {
		return fmt.Errorf("The command '%s' does not exist", cmd.name)
	}

	return handlerFunc(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("the login handler expects a single argument, the username")
	}

	username := cmd.args[0]
	_, err := s.database.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	err = s.config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println("The user has been set to:", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("the register handler expects a single argument, the username")
	}

	username := cmd.args[0]
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	user, err := s.database.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	err = s.config.SetUser(username)
	if err != nil {
		return err
	}

	s.config.CurrentUserName = user.Name
	fmt.Println("The user has been registered")
	s.config.Print()

	return nil
}
