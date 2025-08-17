📌 Jalankan Test

Di root gateway/:

go test ./middleware -v


Output yang diharapkan:

=== RUN   TestJWTMiddleware_WithValidToken
--- PASS: TestJWTMiddleware_WithValidToken (0.00s)
=== RUN   TestJWTMiddleware_NoToken
--- PASS: TestJWTMiddleware_NoToken (0.00s)
=== RUN   TestJWTMiddleware_InvalidHeader
--- PASS: TestJWTMiddleware_InvalidHeader (0.00s)
PASS
ok  	gateway/middleware	0.003s


👉 Dengan begini, kita yakin:

Kalau header valid → token masuk ke context.

Kalau tidak ada / salah format → context tetap aman (nil).

=====================================================
📌 Asumsi Resolver me

Pastikan resolver me kamu kayak gini (atau mirip):

// graph/schema.resolvers.go
package graph

import (
	"context"
	"gateway/middleware"
)

func (r *queryResolver) Me(ctx context.Context) (string, error) {
	token, _ := ctx.Value(middleware.AuthTokenKey).(string)
	if token == "" {
		return "unauthenticated", nil
	}
	return "hello user with token: " + token, nil
}

📌 Jalankan Test
go test ./graph -v


Output expected:

=== RUN   TestQueryMe_WithToken
--- PASS: TestQueryMe_WithToken (0.00s)
=== RUN   TestQueryMe_NoToken
--- PASS: TestQueryMe_NoToken (0.00s)
PASS
ok  	gateway/graph	0.004s
===================================
