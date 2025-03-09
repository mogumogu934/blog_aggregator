package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mogumogu934/blog_aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("usage: login <username>")
	}

	username := cmd.args[0]

	s.config.SetUser(username)
	fmt.Printf("user set to %s\n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("usage: register <name>")
	}

	username := cmd.args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	ctx := context.Background()

	newUser, err := s.db.CreateUser(ctx, params)
	if err != nil {
		fmt.Printf("error creating user: %s, %v\n", username, err)
		os.Exit(1)
	}

	s.config.SetUser(username)
	fmt.Println("User created successfully")
	fmt.Println("User details:")
	fmt.Printf("  ID: %s\n", newUser.ID)
	fmt.Printf("  Name: %s\n", newUser.Name)
	fmt.Printf("  Created At: %s\n", newUser.CreatedAt)
	fmt.Printf("  Updated At: %s\n", newUser.UpdatedAt)

	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		fmt.Printf("error getting users: %v\n", err)
		os.Exit(1)
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
