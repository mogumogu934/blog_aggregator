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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("usage: follow <url>")
	}

	ctx := context.Background()
	url := cmd.args[0]

	feedInfo, err := s.db.GetFeedIDAndNameFromURL(ctx, url)
	if err != nil {
		fmt.Printf("error getting feed ID and name from URL: %s, %v\n", url, err)
		os.Exit(1)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedInfo.ID,
	}

	_, err = s.db.CreateFeedFollow(ctx, params)
	if err != nil {
		fmt.Printf("error following feed: %s\n", url)
		os.Exit(1)
	}

	fmt.Println("Feed followed successfully")
	fmt.Printf("  Feed: %s\n", feedInfo.Name)
	fmt.Printf("  Current User: %s\n", user.Name)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		fmt.Printf("error getting feed follows for current user: %s, %s, %v\n", user.ID, user.Name, err)
		os.Exit(1)
	}

	for _, follow := range follows {
		fmt.Printf("%s\n", follow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("usage: unfollow <url>")
	}

	url := cmd.args[0]

	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    url,
	}

	ctx := context.Background()

	if err := s.db.DeleteFeedFollow(ctx, params); err != nil {
		fmt.Printf("error deleting feed follow record: %v\n", err)
		os.Exit(1)
	}

	return nil
}
