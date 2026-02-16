package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/lib/pq"
	"orm"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer wraps the testcontainer postgres container
type PostgresContainer struct {
	Container *postgres.PostgresContainer
	Client    *orm.Client
	ConnStr   string
}

// NewPostgresContainer creates a new postgres container for testing
func NewPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	// Get path to sql directory (relative to this file)
	_, filename, _, _ := runtime.Caller(0)
	sqlDir := filepath.Join(filepath.Dir(filename), "sql")

	// Find all .sql files in the sql directory (sorted alphabetically)
	initScripts, err := filepath.Glob(filepath.Join(sqlDir, "*.sql"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob sql files: %w", err)
	}

	opts := []testcontainers.ContainerCustomizer{
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithInitScripts(initScripts...),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	}
	if os.Getenv("TESTCONTAINERS_USE_PODMAN") == "1" {
		opts = append(opts, testcontainers.WithProvider(testcontainers.ProviderPodman))
	}
	container, err := postgres.Run(ctx, "postgres:16-alpine", opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	// Create ent client and run schema migrations
	client, err := orm.NewPostgresClient(ctx, connStr)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("failed to create orm client: %w", err)
	}

	return &PostgresContainer{
		Container: container,
		Client:    client,
		ConnStr:   connStr,
	}, nil
}

// Close terminates the container and closes connections
func (p *PostgresContainer) Close(ctx context.Context) error {
	if p.Client != nil {
		_ = p.Client.Close()
	}
	if p.Container != nil {
		return p.Container.Terminate(ctx)
	}
	return nil
}

// TruncateTables truncates specified tables for test isolation
func (p *PostgresContainer) TruncateTables(ctx context.Context, tables ...string) error {
	db, err := sql.Open("postgres", p.ConnStr)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	for _, table := range tables {
		_, err := db.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}
	return nil
}
