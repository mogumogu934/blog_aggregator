package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mogumogu934/blog_aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		currentUser, err := s.db.GetUser(ctx, s.config.CurrentUserName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				fmt.Printf("error: user with name %s does not exist", s.config.CurrentUserName)
			}
			return fmt.Errorf("error checking if current user %s exists: %w", s.config.CurrentUserName, err)
		}
		return handler(s, cmd, currentUser)
	}
}
