package main

import (
	"errors"
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("usage: agg <time between requests>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error parsing time string to time duration value: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for {
		scrapeFeeds(s)
		<-ticker.C
	}
}
