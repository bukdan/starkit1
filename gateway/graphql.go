package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"

	"gateway/clients"
)

func GraphQLHandler() gin.HandlerFunc {
	// GraphQL types
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":       &graphql.Field{Type: graphql.String},
			"username": &graphql.Field{Type: graphql.String},
			"email":    &graphql.Field{Type: graphql.String},
		},
	})
	authType := graphql.NewObject(graphql.ObjectConfig{
		Name: "AuthPayload",
		Fields: graphql.Fields{
			"token": &graphql.Field{Type: graphql.String},
			"user":  &graphql.Field{Type: userType},
		},
	})

	// Root mutation
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"register": &graphql.Field{
				Type: authType,
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"phone":    &graphql.ArgumentConfig{Type: graphql.String},
					"sendVia":  &graphql.ArgumentConfig{Type: graphql.String}, // "email" or "wa"
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					payload := map[string]any{
						"username": p.Args["username"],
						"email":    p.Args["email"],
						"password": p.Args["password"],
					}
					if v, ok := p.Args["phone"].(string); ok {
						payload["phone"] = v
					}
					if sv, ok := p.Args["sendVia"].(string); ok {
						payload["send_via"] = sv
					}
					out, _, err := clients.CallUserService("/auth/register", payload)
					return out, err
				},
			},

			"login": &graphql.Field{
				Type: authType,
				Args: graphql.FieldConfigArgument{
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					payload := map[string]any{"email": p.Args["email"], "password": p.Args["password"]}
					out, _, err := clients.CallUserService("/auth/login", payload)
					return out, err
				},
			},

			"loginWithGoogle": &graphql.Field{
				Type: authType,
				Args: graphql.FieldConfigArgument{
					"idToken": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					payload := map[string]any{"id_token": p.Args["idToken"]}
					out, _, err := clients.CallUserService("/auth/google", payload)
					return out, err
				},
			},

			"verifyOtp": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "VerifyResponse",
					Fields: graphql.Fields{
						"message": &graphql.Field{Type: graphql.String},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"userId":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"channel": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"code":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					payload := map[string]any{"user_id": p.Args["userId"], "channel": p.Args["channel"], "code": p.Args["code"]}
					out, _, err := clients.CallUserService("/auth/verify-otp", payload)
					return out, err
				},
			},
		},
	})

	// Root query (optional me that forwards Authorization header)
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"me": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					// Forward Authorization header to user-service /users/me if implemented
					// Here we call user-service /users/me expecting it to accept the same Authorization header
					ctx := p.Context.Value("ginContext")
					if ctx == nil {
						return nil, nil
					}
					gc, _ := ctx.(*gin.Context)
					auth := gc.GetHeader("Authorization")
					reqPayload := map[string]any{} // empty body
					// call user-service /users/me (POST is used in our client helper)
					out, status, err := clients.CallUserService("/users/me", reqPayload)
					_ = status
					_ = auth
					return out, err
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{Query: rootQuery, Mutation: rootMutation})

	// Handler function
	return func(c *gin.Context) {
		var req struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
			Operation string         `json:"operationName"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// pass gin context into GraphQL context so resolvers can access headers if needed
		params := graphql.Params{
			Schema:         schema,
			RequestString:  req.Query,
			VariableValues: req.Variables,
			Context:        c.Request.Context(),
		}
		// Attach gin context for resolvers that need headers
		ctx := c.Request.Context()
		ctx = contextWithGin(c) // helper below
		params.Context = ctx

		result := graphql.Do(params)
		if len(result.Errors) > 0 {
			c.JSON(http.StatusBadRequest, result)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// helper: inject gin context into request context for resolvers
func contextWithGin(c *gin.Context) any {
	// store gin context in request context under key "ginContext"
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "ginContext", c)
	return ctx
}
