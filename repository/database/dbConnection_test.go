package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestConnectDB(t *testing.T) {
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
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	defer func() {
		if err := ctr.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate PostgreSQL container: %v", err)
		}
	}()

	dbURL, err := ctr.ConnectionString(ctx)
	require.NoError(t, err)

	dbConfig := Config{
		DSN:             dbURL,
		MaxConns:        10,
		MinConns:        1,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}
	// pool, err := ConnectDB(dbConfig)
	// if err != nil {
	// 	t.Fatalf("Failed to connect to database: %v", err)
	// }

	// Create the migrated_template database
	// _, err = pool.Exec(ctx, "CREATE DATABASE migrated_template TEMPLATE test-db")
	// if err != nil {
	// 	t.Fatalf("Failed to create template database: %v", err)
	// }

	// Run any migrations on the database
	_, _, err = ctr.Exec(ctx, []string{"psql", "-U", "postgres", "-d", "test-db", "-c", `CREATE TABLE users (
	
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		sub_id VARCHAR(50),
		verification_status BOOLEAN,
		setup_status VARCHAR(50),
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now()
	  );`})
	require.NoError(t, err)

	// 2. Create a snapshot of the database to restore later
	// tt.options comes the test case, it can be specified as e.g. `postgres.WithSnapshotName("custom-snapshot")` or omitted, to use default name
	err = ctr.Snapshot(ctx)
	require.NoError(t, err)

	/*var currentTime time.Time
	err = pool.QueryRow(ctx, "SELECT NOW()").Scan(&currentTime)
	if err != nil {
		t.Fatalf("Failed to execute test query: %v", err)
	}

	// Log the success message
	t.Logf("Successfully connected to the database at %v", currentTime)*/

	t.Run("Test inserting a user", func(t *testing.T) {
		t.Cleanup(func() {
			// 3. In each test, reset the DB to its snapshot state.
			err = ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pool, err := ConnectDB(dbConfig)
		if err != nil {
			t.Fatalf("Failed to connect to database: %v", err)
		}
		require.NoError(t, err)
		// defer pool.Close(context.Background())
		defer pool.Close()

		_, err = pool.Exec(ctx, `INSERT INTO users (
	  name,
	  email,
	  sub_id,
	  verification_status,
	  setup_status,
	  created_at,
	  updated_at
	  )
	  VALUES (
	  'John Doe',
	   'john.doe@example.com',
	   'SUB987654321',
	   TRUE,
	   'in_progress',
	   NOW(),
	   NOW()
	  )
	   ;`,
		// "John Doe", "john.doe@example.com", "SUB987654321", true, "in_progress", time.Now(), time.Now()
		)
		require.NoError(t, err)
	})
	/*t.Run("Test inserting a user", func(t *testing.T) {
		t.Cleanup(func() {
			// 3. In each test, reset the DB to its snapshot state.
			err = ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pool, err := ConnectDB(dbConfig)
		if err != nil {
			t.Fatalf("Failed to connect to database: %v", err)
		}
		// Ensure the pool is closed after the test
		// defer pool.Close()
		// conn, err := pgx.Connect(context.Background(), dbURL)
		require.NoError(t, err)
		// defer pool.Close(context.Background())
		defer pool.Close()

		_, err = pool.Exec(ctx, "INSERT INTO users(name, age) VALUES ($1, $2)", "test", 42)
		require.NoError(t, err)

		var name string
		var age int64
		err = pool.QueryRow(context.Background(), "SELECT name, age FROM users LIMIT 1").Scan(&name, &age)
		require.NoError(t, err)

		require.Equal(t, "test", name)
		require.EqualValues(t, 42, age)
	})*/
	// Start a PostgreSQL container
	/*pgContainer, err := postgres.Run(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join("..", "testdata", "init-db.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})
	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	dbConfig := Config{
		DSN:             dsn,
		MaxConns:        10,
		MinConns:        1,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}

	pool, err := ConnectDB(dbConfig)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	var currentTime time.Time
	err = pool.QueryRow(ctx, "SELECT NOW()").Scan(&currentTime)
	if err != nil {
		t.Fatalf("Failed to execute test query: %v", err)
	}

	t.Logf("Successfully connected to the database at %v", currentTime)*/
}
