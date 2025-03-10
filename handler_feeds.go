package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mogumogu934/blog_aggregator/internal/database"
	"github.com/mogumogu934/blog_aggregator/internal/fetch"
)

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
		return fmt.Errorf("error creating feed: %w", err)
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
		return fmt.Errorf("error creating feed follow record: %w", err)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()

	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}

	if len(feeds) == 0 {
		return errors.New("you have yet to add any feeds")
	}

	for _, feed := range feeds {
		creatorName, err := s.db.GetUserNameFromID(ctx, feed.UserID)
		if err != nil {
			fmt.Printf("error getting creator name from ID %s: %v\n", feed.UserID, err)
		}

		fmt.Printf("%s\n", feed.Name)
		fmt.Printf("  %s\n", feed.Url)
		fmt.Printf("  %s\n", creatorName)
	}

	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	totalPosts := 0
	newPosts := 0

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch: %w", err)
	}

	if err = s.db.MarkFeedFetched(ctx, nextFeed.ID); err != nil {
		return fmt.Errorf("error marking feed as fetched: %w", err)
	}

	feed, err := fetch.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	fmt.Printf("Feed: %s\n", feed.Channel.Title)

	for _, item := range feed.Channel.Item {
		descriptionSQL := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		pubDateSQL, err := parsePubDate(item.PubDate)
		if err != nil {
			log.Println("error parsing date:", err)
		}

		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: descriptionSQL,
			PublishedAt: pubDateSQL,
			FeedID:      nextFeed.ID,
		}

		_, err = s.db.CreatePost(ctx, params)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique constraint") {
				totalPosts++
				continue
			} else {
				log.Printf("error creating post %s: %v\n", item.Title, err)
			}
		}

		fmt.Printf("  Saved post: %s\n", item.Title)
		totalPosts++
		newPosts++
	}

	fmt.Printf("  Processed %d posts from %s (%d new)\n", totalPosts, feed.Channel.Title, newPosts)
	fmt.Println()

	return nil
}

func parsePubDate(dateStr string) (sql.NullTime, error) {
	if dateStr == "" {
		return sql.NullTime{Valid: false}, nil
	}

	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
		"2006-01-02T15:04:05-07:00",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		parsedTime, err := time.Parse(format, dateStr)
		if err == nil {
			return sql.NullTime{Time: parsedTime, Valid: true}, nil
		}
	}

	return sql.NullTime{Valid: false}, fmt.Errorf("unable to parse date: %s", dateStr)
}
