package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func GraphQLHandler() gin.HandlerFunc {
	// User GraphQL type
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.String},
			"username":    &graphql.Field{Type: graphql.String},
			"email":       &graphql.Field{Type: graphql.String},
			"phone":       &graphql.Field{Type: graphql.String},
			"avatar_url":  &graphql.Field{Type: graphql.String},
			"role":        &graphql.Field{Type: graphql.String},
			"is_verified": &graphql.Field{Type: graphql.Boolean},
		},
	})

	authType := graphql.NewObject(graphql.ObjectConfig{
		Name: "AuthPayload",
		Fields: graphql.Fields{
			"token": &graphql.Field{Type: graphql.String},
			"user":  &graphql.Field{Type: userType},
		},
	})

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
					if v, ok := p.Args["phone"].(string); ok && v != "" {
						payload["phone"] = v
					}
					if sv, ok := p.Args["sendVia"].(string); ok && sv != "" {
						payload["send_via"] = sv
					} // user-service expects send_via
					out, _, err := postJSON("/auth/register", payload, nil)
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
					out, _, err := postJSON("/auth/login", payload, nil)
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
					out, _, err := postJSON("/auth/google", payload, nil)
					return out, err
				},
			},

			"verifyOtp": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "VerifyResp",
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
					out, _, err := postJSON("/auth/verify-otp", payload, nil)
					return out, err
				},
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"me": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					// forward Authorization header if present
					ctx := p.Context
					gc, ok := ctx.Value("gin").(*gin.Context)
					headers := map[string]string{}
					if ok {
						if auth := gc.GetHeader("Authorization"); strings.TrimSpace(auth) != "" {
							headers["Authorization"] = auth
						}
					}
					// call user-service /users/me (GET)
					out, _, err := getJSON("/users/me", headers)
					return out, err
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{Query: rootQuery, Mutation: rootMutation})

	return func(c *gin.Context) {
		var req struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
			OpName    string         `json:"operationName"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// create context that allows resolvers to access gin context
		ctx := context.WithValue(c.Request.Context(), "gin", c)

		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  req.Query,
			VariableValues: req.Variables,
			OperationName:  req.OpName,
			Context:        ctx,
		})

		if len(result.Errors) > 0 {
			c.JSON(http.StatusBadRequest, result)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
