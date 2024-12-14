package tests

import (
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"

	//----------------
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/Wigglor/webservice-v2/repository/database"
	"github.com/Wigglor/webservice-v2/router"
)

func TestConnectDB(t *testing.T) {
	// Initialize the test setup
	setup := SetupTestDB(t)

	t.Run("Test inserting a user", func(t *testing.T) {
		t.Cleanup(func() {
			setup.RestoreTestDB(t)
		})

		pool, err := database.ConnectDB(setup.DBConfig)
		require.NoError(t, err)
		defer pool.Close()

		_, err = pool.Exec(setup.Context, `INSERT INTO users (
			id, name, email, sub_id, verification_status, setup_status, created_at, updated_at
		) VALUES (
			1, 'John Doe', 'john.doe@example.com', 'SUB987654321', TRUE, 'in_progress', NOW(), NOW()
		);`)
		require.NoError(t, err)

		var name string
		var id int32
		err = pool.QueryRow(setup.Context, "SELECT name, id FROM users LIMIT 1").Scan(&name, &id)
		require.NoError(t, err)

		require.Equal(t, "John Doe", name)
		require.EqualValues(t, 1, id)
	})
}

func TestHandlers2(t *testing.T) {
	// Initialize the test setup
	setup := SetupTestDB(t)

	t.Cleanup(func() {
		setup.RestoreTestDB(t)
	})

	// Connect to the database
	pool, err := database.ConnectDB(setup.DBConfig)
	require.NoError(t, err)
	defer pool.Close()

	// Set up the router with the database
	router := router.SetupRouter(pool)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	client := testServer.Client()
	resp, err := client.Get(testServer.URL + "/api/users")
	require.NoError(t, err)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	t.Logf("Response: %s", string(bodyBytes))
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

/*import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/Wigglor/webservice-v2/router"
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

	// 2. Create a snapshot of the database to restore later
	// tt.options comes the test case, it can be specified as e.g. `postgres.WithSnapshotName("custom-snapshot")` or omitted, to use default name
	err = ctr.Snapshot(ctx)
	require.NoError(t, err)

	// var currentTime time.Time
	// err = pool.QueryRow(ctx, "SELECT NOW()").Scan(&currentTime)
	// if err != nil {
	// 	t.Fatalf("Failed to execute test query: %v", err)
	// }

	// // Log the success message
	// t.Logf("Successfully connected to the database at %v", currentTime)

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
		id,
	  name,
	  email,
	  sub_id,
	  verification_status,
	  setup_status,
	  created_at,
	  updated_at
	  )
	  VALUES (
	  1,
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

		var name string
		// var id int64
		var id int32
		err = pool.QueryRow(context.Background(), "SELECT name, id FROM users LIMIT 1").Scan(&name, &id)
		require.NoError(t, err)

		require.Equal(t, "John Doe", name)
		require.EqualValues(t, 1, id)
		// -----------------------------------------------
		// -----------------------------------------------
		router := router.SetupRouter(pool)

		testServer := httptest.NewServer(router)
		println(testServer.URL)
		println(testServer.URL + "/api/users")
		defer testServer.Close()
		client := testServer.Client()
		req, err := http.NewRequest("GET", testServer.URL+"/api/users", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		// t.Logf("Response body: %s", resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		t.Logf("Response body: %s", string(bodyBytes))
		// -----------------------------------------------
		// -----------------------------------------------
		userUrl := testServer.URL + "/api/user/1"
		req2, err := http.NewRequest("GET", userUrl, nil)
		// req2, err := http.NewRequest("GET", fmt.Sprintf("%s/api/user/%d", testServer.URL, 1), nil)
		// 	require.NoError(t, err)
		require.NoError(t, err)
		req2.Header.Set("Content-Type", "application/json")

		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()

		bodyBytes2, err := io.ReadAll(resp2.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		t.Logf("URL string: %s", userUrl)
		t.Logf("Response body user ID: %s", string(bodyBytes2))

		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Test getting all users", func(t *testing.T) {
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

		_, err = pool.Exec(ctx, "SELECT * FROM users") // "John Doe", "john.doe@example.com", "SUB987654321", true, "in_progress", time.Now(), time.Now()
		require.NoError(t, err)
	})

	t.Run("Test get all user by id - httptest", func(t *testing.T) {
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
			id,
		  name,
		  email,
		  sub_id,
		  verification_status,
		  setup_status,
		  created_at,
		  updated_at
		  )
		  VALUES (
		  1,
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

		router := router.SetupRouter(pool)

		testServer := httptest.NewServer(router)
		defer testServer.Close()
		client := testServer.Client()

		userUrl := testServer.URL + "/api/user/1"
		req, err := http.NewRequest("GET", userUrl, nil)
		// req2, err := http.NewRequest("GET", fmt.Sprintf("%s/api/user/%d", testServer.URL, 1), nil)
		// 	require.NoError(t, err)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		bodyBytes2, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		t.Logf("URL string: %s", userUrl)
		t.Logf("Response body user ID: %s", string(bodyBytes2))

		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// t.Run("Test getting all users - httptest", func(t *testing.T) {
	// 	t.Cleanup(func() {
	// 		// 3. In each test, reset the DB to its snapshot state.
	// 		err = ctr.Restore(ctx)
	// 		require.NoError(t, err)
	// 	})

	// 	pool, err := ConnectDB(dbConfig)
	// 	if err != nil {
	// 		t.Fatalf("Failed to connect to database: %v", err)
	// 	}
	// 	require.NoError(t, err)
	// 	// defer pool.Close(context.Background())
	// 	defer pool.Close()

	// 	router := router.SetupRouter(pool)

	// 	testServer := httptest.NewServer(router)
	// 	println(testServer.URL)
	// 	println(testServer.URL + "/api/users")
	// 	defer testServer.Close()
	// 	client := testServer.Client()
	// 	req, err := http.NewRequest("GET", testServer.URL+"/api/users", nil)
	// 	require.NoError(t, err)
	// 	req.Header.Set("Content-Type", "application/json")

	// 	resp, err := client.Do(req)
	// 	// t.Logf("Response body: %s", resp.Body)
	// 	require.NoError(t, err)
	// 	defer resp.Body.Close()

	// 	bodyBytes, err := io.ReadAll(resp.Body)
	// 	if err != nil {
	// 		t.Fatalf("Failed to read response body: %v", err)
	// 	}
	// 	t.Logf("Response body: %s", string(bodyBytes))

	// 	require.Equal(t, http.StatusOK, resp.StatusCode)
	// })

	// t.Run("Test get all user by id - httptest", func(t *testing.T) {
	// 	t.Cleanup(func() {
	// 		// 3. In each test, reset the DB to its snapshot state.
	// 		err = ctr.Restore(ctx)
	// 		require.NoError(t, err)
	// 	})

	// 	pool, err := ConnectDB(dbConfig)
	// 	if err != nil {
	// 		t.Fatalf("Failed to connect to database: %v", err)
	// 	}
	// 	require.NoError(t, err)
	// 	// defer pool.Close(context.Background())
	// 	defer pool.Close()

	// 	router := router.SetupRouter(pool)

	// 	testServer := httptest.NewServer(router)

	// 	defer testServer.Close()
	// 	client := testServer.Client()
	// 	req, err := http.NewRequest("GET", testServer.URL+"/api/user/1", nil)
	// 	require.NoError(t, err)
	// 	req.Header.Set("Content-Type", "application/json")

	// 	resp, err := client.Do(req)
	// 	require.NoError(t, err)
	// 	defer resp.Body.Close()

	// 	require.Equal(t, http.StatusOK, resp.StatusCode)
	// })
	// t.Run("Test inserting a user", func(t *testing.T) {
	// 	t.Cleanup(func() {
	// 		// 3. In each test, reset the DB to its snapshot state.
	// 		err = ctr.Restore(ctx)
	// 		require.NoError(t, err)
	// 	})

	// 	pool, err := ConnectDB(dbConfig)
	// 	if err != nil {
	// 		t.Fatalf("Failed to connect to database: %v", err)
	// 	}
	// 	// Ensure the pool is closed after the test
	// 	// defer pool.Close()
	// 	// conn, err := pgx.Connect(context.Background(), dbURL)
	// 	require.NoError(t, err)
	// 	// defer pool.Close(context.Background())
	// 	defer pool.Close()

	// 	_, err = pool.Exec(ctx, "INSERT INTO users(name, age) VALUES ($1, $2)", "test", 42)
	// 	require.NoError(t, err)

	// 	var name string
	// 	var age int64
	// 	err = pool.QueryRow(context.Background(), "SELECT name, age FROM users LIMIT 1").Scan(&name, &age)
	// 	require.NoError(t, err)

	// 	require.Equal(t, "test", name)
	// 	require.EqualValues(t, 42, age)
	// })
	// Start a PostgreSQL container
	//---------------------
	// pgContainer, err := postgres.Run(ctx,
	// 	testcontainers.WithImage("postgres:15.3-alpine"),
	// 	postgres.WithInitScripts(filepath.Join("..", "testdata", "init-db.sql")),
	// 	postgres.WithDatabase("test-db"),
	// 	postgres.WithUsername("postgres"),
	// 	postgres.WithPassword("postgres"),
	// 	testcontainers.WithWaitStrategy(
	// 		wait.ForLog("database system is ready to accept connections").
	// 			WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	// )
	// if err != nil {
	// 	t.Fatalf("Failed to start PostgreSQL container: %v", err)
	// }

	// t.Cleanup(func() {
	// 	if err := pgContainer.Terminate(ctx); err != nil {
	// 		t.Fatalf("failed to terminate pgContainer: %s", err)
	// 	}
	// })
	// dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	// assert.NoError(t, err)

	// dbConfig := Config{
	// 	DSN:             dsn,
	// 	MaxConns:        10,
	// 	MinConns:        1,
	// 	MaxConnLifetime: time.Hour,
	// 	MaxConnIdleTime: 30 * time.Minute,
	// }

	// pool, err := ConnectDB(dbConfig)
	// if err != nil {
	// 	t.Fatalf("Failed to connect to database: %v", err)
	// }
	// defer pool.Close()

	// var currentTime time.Time
	// err = pool.QueryRow(ctx, "SELECT NOW()").Scan(&currentTime)
	// if err != nil {
	// 	t.Fatalf("Failed to execute test query: %v", err)
	// }

	// t.Logf("Successfully connected to the database at %v", currentTime)
}
*/
