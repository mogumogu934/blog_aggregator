package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mogumogu934/blog_aggregator/internal/database"
)

func cleanDescription(description string) string {
	policy := bluemonday.StripTagsPolicy()
	return policy.Sanitize(description)
}

func truncateString(description string, maxLen int) string {
	runes := []rune(description)

	if len(runes) <= maxLen {
		return description
	}

	return string(runes[:maxLen-3]) + "..."
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var postLimit int32 = 2 // Default value
	if len(cmd.args) > 0 {
		postLimitInt64, err := strconv.ParseInt(cmd.args[0], 10, 32) // Convert string argument to int. ParseInt always returns int64.
		if err != nil {
			return fmt.Errorf("error parsing limit argument %s to int64: %w", cmd.args[0], err)
		}

		postLimit = int32(postLimitInt64)
	}

	ctx := context.Background()
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  postLimit,
	}

	posts, err := s.db.GetPostsForUser(ctx, params)
	if err != nil {
		return fmt.Errorf("error getting posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found. Try following some feeds first!")
		return nil
	}

	for _, post := range posts {
		dateStr := "Unknown date"
		if post.PublishedAt.Valid {
			dateStr = post.PublishedAt.Time.Format("Jan 02, 2006")
		}

		fmt.Printf("%s\n", post.Title)
		fmt.Printf("  %s, %s\n", dateStr, post.Url)

		/*
			if post.Description.Valid {
				fmt.Printf("DEBUG RAW DESCRIPTION: %s\n", post.Description.String)
			}
		*/

		if post.Description.Valid && post.Description.String != "" {
			cleanDesc := cleanDescription(post.Description.String)
			truncatedDesc := truncateString(cleanDesc, 203)
			fmt.Printf("  %s\n", truncatedDesc)
		}
		fmt.Println()
	}

	return nil
}
