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
	"github.com/mogumogu934/blog_aggregator/internal/fetch"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("usage: agg <time between requests>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		fmt.Printf("error parsing time string to time duration value: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("usage: addfeed <name> <url>")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	addFeedParams := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Valid: false,
			Time:  time.Time{},
		},
		Name:   name,
		Url:    url,
		UserID: user.ID,
	}

	ctx := context.Background()

	feed, err := s.db.AddFeed(ctx, addFeedParams)
	if err != nil {
		fmt.Printf("error creating feed: %v\n", err)
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

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(ctx, feedFollowParams)
	if err != nil {
		fmt.Printf("error creating feed follow record: %v\n", err)
		os.Exit(1)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()

	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		fmt.Printf("error getting feeds: %v", err)
		os.Exit(1)
	}

	for _, feed := range feeds {
		creatorName, err := s.db.GetUserNameFromID(ctx, feed.UserID)
		if err != nil {
			fmt.Printf("error getting creator name from ID: %s, %v\n", feed.UserID, err)
		}

		fmt.Printf("%s\n", feed.Name)
		fmt.Printf("  %s\n", feed.Url)
		fmt.Printf("  %s\n", creatorName)
	}

	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		fmt.Printf("error getting next feed to fetch: %v\n", err)
		os.Exit(1)
	}

	if err = s.db.MarkFeedFetched(ctx, nextFeed.ID); err != nil {
		fmt.Printf("error marking feed as fetched: %v\n", err)
		os.Exit(1)
	}

	feed, err := fetch.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		fmt.Printf("error fetching feed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Feed: %s\n", feed.Channel.Title)
	for _, item := range feed.Channel.Item {
		fmt.Printf("  %s\n", item.Title)
	}
	fmt.Println()

	return nil
}
