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

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return nil
	}

	fmt.Printf("feed: %v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("the register handler expects two arguments, the feed name and the url")
	}

	feedName := cmd.args[0]
	url := cmd.args[1]
	user, err := s.database.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       url,
		UserID:    user.ID,
	}
	_, err = s.database.CreateFeed(context.Background(), feed)
	return err

}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.database.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Feeds")
	fmt.Println("---------------")
	for _, feed := range feeds {
		fmt.Println("Name:", feed.Name)
		fmt.Println("URL:", feed.Url)
		fmt.Println("User:", feed.Username)
	}
	fmt.Println()
	return nil
}
