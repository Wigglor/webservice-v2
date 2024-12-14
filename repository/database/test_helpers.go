package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TestSetup struct {
	Container *postgres.PostgresContainer
	DBConfig  Config
	Context   context.Context
}

// SetupTestDB initializes the test container, config, and context
func SetupTestDB(t *testing.T) *TestSetup {
	ctx := context.Background()

	ctr, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	require.NoError(t, err)

	// Ensure the container is cleaned up after the test
	t.Cleanup(func() {
		if err := ctr.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate PostgreSQL container: %v", err)
		}
	})

	dbURL, err := ctr.ConnectionString(ctx)
	require.NoError(t, err)

	dbConfig := Config{
		DSN:             dbURL,
		MaxConns:        10,
		MinConns:        1,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}

	_, _, err = ctr.Exec(ctx, []string{"psql", "-U", "postgres", "-d", "test-db", "-c", `CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		sub_id VARCHAR(50),
		verification_status BOOLEAN,
		setup_status VARCHAR(50),
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now()
	);`})
	require.NoError(t, err)

	// Create a snapshot for restoring
	err = ctr.Snapshot(ctx)
	require.NoError(t, err)

	return &TestSetup{
		Container: ctr,
		DBConfig:  dbConfig,
		Context:   ctx,
	}
}

// RestoreTestDB restores the database to its initial state
func (ts *TestSetup) RestoreTestDB(t *testing.T) {
	err := ts.Container.Restore(ts.Context)
	require.NoError(t, err)
}
