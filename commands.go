package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mogumogu934/blog_aggregator/internal/config"
	"github.com/mogumogu934/blog_aggregator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("usage: login <username>")
	}

	username := cmd.args[0]
	ctx := context.Background()

	_, err := s.db.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("error: user with name %s does not exist\n", username)
			os.Exit(1)
		}
		return fmt.Errorf("error checking if user exists: %w", err)
	}

	s.config.SetUser(username)
	fmt.Printf("user set to %s\n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("not enough arguments were provided")
	}

	username := cmd.args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	ctx := context.Background()

	_, err := s.db.GetUser(ctx, username)
	if err == nil {
		fmt.Printf("error: user with name %s already exists", username)
		os.Exit(1)
	} else if !errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("error checking if user exists: %v\n", err)
		os.Exit(1)
	}

	user, err := s.db.CreateUser(ctx, params)
	if err != nil {
		fmt.Printf("error creating user: %v: %v\n", params, err)
		os.Exit(1)
	}

	s.config.SetUser(username)
	fmt.Println("User created successfully")
	fmt.Println("User details:")
	fmt.Printf("  ID: %s\n", user.ID)
	fmt.Printf("  Name: %s\n", user.Name)
	fmt.Printf("  Created At: %s\n", user.CreatedAt)
	fmt.Printf("  Updated At: %s\n", user.UpdatedAt)

	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		fmt.Printf("error getting users: %v\n", err)
	}

	for _, user := range users {
		if user == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	if err := s.db.DeleteAllUsers(ctx); err != nil {
		fmt.Printf("error resetting database: %v", err)
		os.Exit(1)
	}
	fmt.Println("Database reset successfully")

	return nil
}

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(ctx, feedURL)
	if err != nil {
		fmt.Printf("error fetching feed: %v", err)
		os.Exit(1)
	}
	fmt.Println(feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("usage: addfeed <name> <url>")
	}

	ctx := context.Background()

	currentUser, err := s.db.GetUser(ctx, s.config.CurrentUserName)
	if err != nil {
		fmt.Printf("error getting current user: %v", err)
		os.Exit(1)
	}

	name := cmd.args[0]
	url := cmd.args[1]

	params := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    currentUser.ID,
	}

	feed, err := s.db.AddFeed(ctx, params)
	if err != nil {
		fmt.Printf("error creating feed: %v: %v\n", params, feed)
		os.Exit(1)
	}

	fmt.Println("Feed created successfully")
	fmt.Println("Feed details:")
	fmt.Printf("  ID: %s\n", feed.ID)
	fmt.Printf("  Created At: %s\n", feed.CreatedAt)
	fmt.Printf("  Updated At: %s\n", feed.UpdatedAt)
	fmt.Printf("  Name: %s\n", feed.Name)
	fmt.Printf("  URL: %s\n", feed.Url)
	fmt.Printf("  User ID: %s\n", feed.UserID)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		fmt.Printf("error getting feeds: %v", err)
	}

	for _, feed := range feeds {
		creatorName, err := s.db.GetUserNameFromID(ctx, feed.UserID)
		if err != nil {
			fmt.Printf("error getting creator name from ID %s: %v\n", feed.UserID, err)
		}

		fmt.Printf("Feed: %s\n", feed.Name)
		fmt.Printf("  URL: %s\n", feed.Url)
		fmt.Printf("  Creator: %s\n", creatorName)
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if handler, exists := c.handlers[cmd.name]; exists {
		return handler(s, cmd)
	}

	return fmt.Errorf("command not found: %s", cmd.name)
}
