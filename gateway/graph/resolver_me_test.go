package graph_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"

	"gateway/graph"
	"gateway/graph/model"
	"gateway/middleware"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
)

func TestQueryMe_WithToken(t *testing.T) {
	// setup server gqlgen dengan resolver
	resolver := &graph.Resolver{}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	c := client.New(srv)

	// bikin context dengan token dummy
	ctx := context.WithValue(context.Background(), middleware.AuthTokenKey, "dummy-token")

	// jalankan query `me`
	var resp struct {
		Me string
	}
	c.MustPost(`
		query {
			me
		}
	`, &resp, client.WithContext(ctx))

	if !strings.Contains(resp.Me, "dummy-token") {
		t.Errorf("expected response contain dummy-token, got %v", resp.Me)
	}
}

func TestQueryMe_NoToken(t *testing.T) {
	resolver := &graph.Resolver{}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	c := client.New(srv)

	var resp struct {
		Me string
	}
	c.MustPost(`
		query {
			me
		}
	`, &resp)

	if resp.Me != "unauthenticated" {
		t.Errorf("expected unauthenticated, got %v", resp.Me)
	}
}

// gateway/graph/schema.resolvers.go
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	token, ok := ctx.Value(middleware.AuthTokenKey).(string)
	if !ok || token == "" {
		return nil, errors.New("unauthenticated")
	}

	req, _ := http.NewRequest("GET", os.Getenv("USER_SERVICE_URL")+"/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user profile")
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
