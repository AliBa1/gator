package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AliBa1/gator/internal/database"
	"github.com/google/uuid"
)

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

func handlerReset(s *state, cmd command) error {
	err := s.database.DeleteUsers(context.Background())
	return err
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.database.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if s.config.CurrentUserName == user.Name {
			fmt.Printf(" (current)")
		}
		fmt.Printf("\n")
	}

	return nil
}
