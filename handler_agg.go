package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("usage: agg <time between requests>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		fmt.Println("error parsing time string to time duration value:", err)
		os.Exit(1)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for {
		scrapeFeeds(s)
		<-ticker.C
	}
}
