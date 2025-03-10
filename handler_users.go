package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
			return fmt.Errorf("error: user with name %s does not exist", username)
		}
		return fmt.Errorf("error checking if user with name %s already exists: %w", username, err)
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
		return fmt.Errorf("error: user with name %s already exists", username)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking if user with name %s already exists: %w", username, err)
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	newUser, err := s.db.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("error creating user %s: %w", username, err)
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
		return fmt.Errorf("error getting users: %w", err)
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
		return fmt.Errorf("error resetting database: %w", err)
	}
	fmt.Println("Database reset successfully")

	return nil
}
