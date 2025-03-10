package main

import (
	"context"
	"errors"
	"fmt"
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

	feedInfo, err := s.db.GetFeedFromURL(ctx, url)
	if err != nil {
		return fmt.Errorf("error getting feed ID and name from URL %s: %v", url, err)
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
		return fmt.Errorf("error following feed %s: %v", url, err)
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
		return fmt.Errorf("error getting feed follows for current user %s: %v", user.Name, err)
	}

	if len(follows) == 0 {
		return errors.New("you have yet to follow any feeds")
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
	ctx := context.Background()

	feed, err := s.db.GetFeedFromURL(ctx, url)
	if err != nil {
		return fmt.Errorf("error getting feed: %w", err)
	}

	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if err := s.db.DeleteFeedFollow(ctx, params); err != nil {
		fmt.Println("error deleting feed follow record:", err)
	}

	return nil
}
