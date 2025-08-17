package main

import (
	"log"
	"net/http"
	"os"

	"gateway/graph"
	"gateway/graph/generated"
	"gateway/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", middleware.JWTMiddleware(srv)) // 🔥 pasang JWT middleware

	log.Printf("🚀 GraphQL Gateway running at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
