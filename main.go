package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/mogumogu934/blog_aggregator/internal/config"
	"github.com/mogumogu934/blog_aggregator/internal/database"
)

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
	commandRegistry.register("reset", handlerReset)

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
