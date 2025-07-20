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
	if len(cmd.args) < 1 {
		return errors.New("the agg handler expects one argument, time between requests (ex: 1s, 1m, 1h)")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("Collecting feeds every", timeBetweenRequests.String())
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("the add feed handler expects two arguments, the feed name and the url")
	}

	feedName := cmd.args[0]
	url := cmd.args[1]

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       url,
		UserID:    user.ID,
	}
	feed, err := s.database.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.database.CreateFeedFollow(context.Background(), feedFollowParams)
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("the follow handler expects one argument, the url")
	}

	url := cmd.args[0]
	feed, err := s.database.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.database.CreateFeedFollow(context.Background(), feedFollowParams)
	return err
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.database.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Println("Feeds Your Following")
	fmt.Println("---------------")
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}
	fmt.Println()
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("the unfollow handler expects one argument, the url")
	}

	url := cmd.args[0]

	deleteFeedFollowParams := database.DeleteFeedFollowParams{
		Url: url,
		ID:  user.ID,
	}
	err := s.database.DeleteFeedFollow(context.Background(), deleteFeedFollowParams)
	if err != nil {
		return err
	}

	return nil
}
