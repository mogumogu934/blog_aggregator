package main

import (
	"fmt"
	"os"

	"github.com/mogumogu934/blog_aggregator/internal/config"
)

func main() {
	appConfig, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config: %v", err)
		os.Exit(1)
	}

	appState := state{
		config: &appConfig,
	}

	commandRegistry := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	commandRegistry.register("login", handlerLogin)

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

	err = commandRegistry.run(&appState, cmd)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
