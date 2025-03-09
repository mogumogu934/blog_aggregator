package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/mogumogu934/blog_aggregator/internal/config"
	"github.com/mogumogu934/blog_aggregator/internal/database"
)

var client *Client

func init() {
	client = NewClient(5 * time.Second)
}

func main() {
	appConfig, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config: %v", err)
		os.Exit(1)
	}

	dbURL := appConfig.DbURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	appState := state{
		db:     dbQueries,
		config: &appConfig,
	}

	commandRegistry := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	commandRegistry.register("login", handlerLogin)
	commandRegistry.register("register", handlerRegister)
	commandRegistry.register("users", handlerUsers)
	commandRegistry.register("reset", handlerReset)
	commandRegistry.register("agg", handlerAgg)
	commandRegistry.register("addfeed", handlerAddFeed)
	commandRegistry.register("feeds", handlerFeeds)
	commandRegistry.register("follow", handlerFollow)
	commandRegistry.register("following", handlerFollowing)

	if len(os.Args) < 2 {
		fmt.Println("error: not enough arguments provided")
		os.Exit(1)
	}

	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	cmd := command{
		name: commandName,
		args: commandArgs,
	}

	if err = commandRegistry.run(&appState, cmd); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
