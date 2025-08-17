package graph_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gateway/graph"
	"gateway/graph/middleware"
	"gateway/graph/model"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
)

// mock user-service server
func startMockUserService(t *testing.T) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer dummy-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user := model.User{
			ID:         "user-123",
			Name:       "Test User",
			Email:      "test@example.com",
			Phone:      "08123456789",
			IsVerified: true,
			Role:       "user",
		}
		json.NewEncoder(w).Encode(user)
	})
	return httptest.NewServer(handler)
}

func TestResolverMe_WithMockUserService(t *testing.T) {
	// start mock server
	mockServer := startMockUserService(t)
	defer mockServer.Close()

	// set USER_SERVICE_URL environment
	old := graph.UserServiceURL
	graph.UserServiceURL = mockServer.URL
	defer func() { graph.UserServiceURL = old }()

	// setup resolver + gql server
	resolver := &graph.Resolver{}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	c := client.New(srv)

	ctx := context.WithValue(context.Background(), middleware.AuthTokenKey, "dummy-token")

	var resp struct {
		Me struct {
			ID    string
			Name  string
			Email string
		}
	}

	c.MustPost(`query { me { id name email } }`, &resp, client.WithContext(ctx))

	if resp.Me.ID != "user-123" {
		t.Errorf("expected user-123, got %s", resp.Me.ID)
	}
	if resp.Me.Name != "Test User" {
		t.Errorf("expected Test User, got %s", resp.Me.Name)
	}
}
