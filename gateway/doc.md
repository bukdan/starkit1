
gateway/
├── go.mod                  # Module Go
├── go.sum                  # Dependencies Go
├── main.go                 # Entry point server GraphQL
├── server.go               # Setup server, register handler & middleware
├── .env.example            # Contoh environment variables
├── config/
│   └── config.go           # Load config/env (USER_SERVICE_URL, JWT secret, PORT)
├── middleware/
│   └── jwt.go              # JWT middleware, inject token ke context
├── graph/
│   ├── schema.graphqls     # Definisi GraphQL schema (Query, Mutation, types)
│   ├── schema.resolvers.go # Implementasi resolver GraphQL (me, login, register)
│   └── model/              # Struct model GraphQL (User, AuthPayload, input types)
├── clients/
│   └── user_client.go      # REST client helper ke user-service (fetch profile, login, register)
├── utils/
│   └── http_client.go      # HTTP helper untuk request antar service
└── Dockerfile              # Dockerfile untuk build image gateway
========================================
🔹 Penjelasan fungsi tiap folder/file
main.go → memulai server GraphQL
server.go → setup HTTP server + apply JWT middleware
.env.example → environment variable (USER_SERVICE_URL, PORT, JWT_SECRET)
config/ → baca env / config global
middleware/jwt.go → ambil JWT dari header, inject ke context
graph/schema.graphqls → definisi Query / Mutation / Types
graph/schema.resolvers.go → implementasi resolver, forward ke user-service
graph/model/ → struct input/output GraphQL (User, AuthPayload)
clients/user_client.go → helper panggil REST API user-service
utils/http_client.go → wrapper HTTP request, handle errors / timeout
Dockerfile → build image untuk deploy Gateway
Cara pakai / run lokal

Buat folder gateway/ dan simpan file-file di atas.

Set env USER_SERVICE_URL kalau user-service ada di container (contoh http://user-service:8081) atau http://localhost:8081 kalau di lokal.

cd gateway && go mod tidy

go run .
GraphQL endpoint: POST http://localhost:8080/graphql

Contoh GraphQL mutations (curl)

Register:

curl -X POST http://localhost:8080/graphql -H "Content-Type: application/json" -d '{
  "query":"mutation ($u:String!, $e:String!, $p:String!) { register(username:$u, email:$e, password:$p, sendVia:\"email\") { token user { id username email } } }",
  "variables": { "u":"alice", "e":"alice@example.com", "p":"secret123" }
}'


Login:

curl -X POST http://localhost:8080/graphql -H "Content-Type: application/json" -d '{
  "query":"mutation { login(email:\"alice@example.com\", password:\"secret123\") { token user { id username email } } }"
}'


Verify OTP:

curl -X POST http://localhost:8080/graphql -H "Content-Type: application/json" -d '{
  "query":"mutation { verifyOtp(userId:\"<user-id>\", channel:\"email\", code:\"123456\") { message } }"
}'


Query me (forward Authorization header):

curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT>" \
  -d '{ "query":"{ me { id username email is_verified } }" }'


Catatan: me resolver forwards Authorization header to user-service. Pastikan user-service implements /users/me (GET) and reads Authorization header.

Catatan & tip integrasi

Field naming: gateway maps sendVia → send_via when calling user-service. Pastikan user-service expects send_via.

Error mapping: currently gateway returns raw JSON from user-service. Kamu bisa map responses to strict GraphQL shapes if needed.

For production, consider:

Use gqlgen for type-safe schema & resolvers.

Add retries / circuit-breaker when calling downstream services.

Add timeout / rate-limiter / auth introspection endpoint.

Use HTTPS between services and mutual auth when needed.
1. Register User
curl -X POST http://localhost:8081/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Budi",
    "email": "budi@example.com",
    "phone": "08123456789",
    "password": "rahasia123"
  }'


✅ Respon (status 201):

{
  "message": "User registered. Please verify OTP."
}

2. Login User
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "budi@example.com",
    "password": "rahasia123"
  }'


✅ Respon (status 200):

{
  "token": "JWT_TOKEN_HERE"
}

3. Verifikasi OTP
curl -X POST http://localhost:8081/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "budi@example.com",
    "code": "123456"
  }'


✅ Respon:

{
  "message": "Account verified successfully."
}

4. Get Profil User (JWT Protected)
curl -X GET http://localhost:8081/users/me \
  -H "Authorization: Bearer JWT_TOKEN_HERE"


✅ Respon:

{
  "id": "uuid-user",
  "name": "Budi",
  "email": "budi@example.com",
  "phone": "08123456789",
  "is_verified": true,
  "role": "user"
}

5. Login dengan Google ID (simulasi)
curl -X POST http://localhost:8081/auth/google \
  -H "Content-Type: application/json" \
  -d '{
    "google_id": "1234567890",
    "email": "budi.google@example.com",
    "name": "Budi Google"
  }'


✅ Respon:

{
  "token": "JWT_TOKEN_FOR_GOOGLE_USER"
}

Test Query di Playground

URL: http://localhost:8080

Register
mutation {
  register(name:"Budi", email:"budi@example.com", password:"rahasia123", phone:"08123456789")
}

Login
mutation {
  login(email:"budi@example.com", password:"rahasia123") {
    token
  }
}
