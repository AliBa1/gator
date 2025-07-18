package config

import (
	"errors"
	"fmt"
)

type State struct {
	Config *Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	commands map[string]func(*State, Command) error
}

func NewCommands() Commands {
	commands := Commands{
		commands: make(map[string]func(*State, Command) error),
	}
	commands.register("login", handlerLogin)
	return commands
}

// func CommandsMap() map[string]func(*State, Command) error {
// 	commands := make(map[string]func(*State, Command) error)
// 	commands["login"] = handlerLogin
// 	return commands
// }

func (c *Commands) Run(s *State, cmd Command) error {
	handlerFunc, ok := c.commands[cmd.Name]
	if !ok {
		return fmt.Errorf("The command '%s' does not exist", cmd.Name)
	}

	return handlerFunc(s, cmd)
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	c.commands[name] = f
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("the login handler expects a single argument, the username")
	}

	username := cmd.Args[0]
	err := s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println("The user has been set to:", username)
	return nil
}
