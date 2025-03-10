package main

import (
	"context"
	"database/sql"
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

	ctx := context.Background()
	username := cmd.args[0]

	_, err := s.db.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("error: user with name %s does not exist\n", username)
			os.Exit(1)
		}
		fmt.Printf("error checking if user with name %s already exists: %v\n", username, err)
		os.Exit(1)
	}

	s.config.SetUser(username)
	fmt.Printf("Current user set to %s\n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("usage: register <name>")
	}

	ctx := context.Background()
	username := cmd.args[0]

	_, err := s.db.GetUser(ctx, username)
	if err == nil {
		fmt.Printf("error: user with name %s already exists\n", username)
		os.Exit(1)
	} else if !errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("error checking if user with name %s already exists: %v\n", username, err)
		os.Exit(1)
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	newUser, err := s.db.CreateUser(ctx, params)
	if err != nil {
		fmt.Printf("error creating user %s: %v\n", username, err)
		os.Exit(1)
	}

	s.config.SetUser(username)
	fmt.Println("User created successfully")
	fmt.Printf("Current user set to %s\n", username)
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
		fmt.Println("error getting users:", err)
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
		fmt.Println("error resetting database:", err)
		os.Exit(1)
	}
	fmt.Println("Database reset successfully")

	return nil
}
