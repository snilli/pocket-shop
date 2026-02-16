package orm

import (
	"context"
	"fmt"

	"orm/ent"

	_ "github.com/lib/pq"
)

// Client wraps the Ent client
type Client struct {
	*ent.Client
}

// NewPostgresClient creates a new ORM client with PostgreSQL connection
func NewPostgresClient(ctx context.Context, dsn string) (*Client, error) {
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	// Run auto migration
	if err := client.Schema.Create(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed creating schema: %w", err)
	}

	return &Client{Client: client}, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	return c.Client.Close()
}
