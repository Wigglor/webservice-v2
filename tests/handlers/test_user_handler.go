package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wigglor/webservice-v2/repository/database"
	"github.com/Wigglor/webservice-v2/router"
	"github.com/Wigglor/webservice-v2/tests"
	"github.com/stretchr/testify/require"
)

func TestHandlers(t *testing.T) {
	// Initialize the test setup
	setup := tests.SetupTestDB(t)

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
