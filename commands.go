package main

import (
	"fmt"

	"github.com/mogumogu934/blog_aggregator/internal/config"
	"github.com/mogumogu934/blog_aggregator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if handler, exists := c.handlers[cmd.name]; exists {
		return handler(s, cmd)
	}

	return fmt.Errorf("command not found: %s", cmd.name)
}
