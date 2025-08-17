
User Service (full)
- Register (with OTP via email/WA)
- Verify OTP
- Login (JWT)
- Profile endpoints
- Admin user management
How to run:
1. Fill .env.example -> .env
2. Ensure PostgreSQL running and DATABASE_URL points to a database 'usersvc'
3. Run migrations: psql -d usersvc -f migrations/001_create_users_and_otp.sql
4. go mod tidy
5. go run .
//////////////////////////////////==============
mau
ChatGPT said:

Oke 🔥 kita langsung siapkan user-service full runnable ZIP.
Isinya sudah termasuk:

Struktur folder user-service/ sesuai yang tadi.

Semua file Go (handler, service, repository, JWT, OTP, hash, Google login).

File migrasi SQL (users, otp_codes).

Dockerfile dan .env.example.

README.md untuk cara jalanin dengan Docker.

Tapi di sini aku nggak bisa langsung kasih file ZIP, jadi aku akan:

Generate semua file lengkap di dalam satu blok kode.

Kamu tinggal copy ke folder user-service/ sesuai struktur.

Lalu zip manual (zip -r user-service-full.zip user-service/).

=====================================================
Step 12: Uji Coba dengan Postman / cURL
1. Register User

👉 Request ke user-service (REST API).

POST http://localhost:8081/register

{
  "email": "test@example.com",
  "password": "secret123",
  "username": "testuser"
}


📌 Expected Response:

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
}


Token ini JWT yang bisa dipakai untuk autentikasi di Gateway.

2. Login User

👉 Request ke user-service (REST API).

POST http://localhost:8081/login

{
  "email": "test@example.com",
  "password": "secret123"
}


📌 Expected Response:

{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}


Simpan token ini, karena dipakai di step berikutnya.

3. Query GraphQL (me) lewat Gateway

👉 Request ke gateway (GraphQL endpoint).

POST http://localhost:8080/query
Header:

Content-Type: application/json
Authorization: Bearer <JWT_TOKEN_DARI_LOGIN>


Body:

{
  "query": "query { me { id email username } }"
}


📌 Expected Response:

{
  "data": {
    "me": {
      "id": "c0b1b3b8-91a7-4e21-8d2f-xxxx",
      "email": "test@example.com",
      "username": "testuser"
    }
  }
}

4. Kalau tanpa token / token invalid

Query ke me akan balas error:

{
  "errors": [
    {
      "message": "unauthorized",
      "path": ["me"]
    }
  ],
  "data": null
}


Jadi flow-nya sekarang:

Register / Login di user-service → dapet JWT.

Kirim JWT ke gateway GraphQL di header Authorization.

Query me akan jalan dan validasi JWT di middleware.
============================
Buat file bernama user-service-test.postman_collection.json:

{
  "info": {
    "name": "User Service Test",
    "_postman_id": "f58e53a5-4e4f-4d3e-a7a0-123456789abc",
    "description": "Test collection for User Service + Gateway GraphQL",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Register User",
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"email\": \"test@example.com\",\n  \"password\": \"secret123\",\n  \"username\": \"testuser\"\n}"
        },
        "url": {
          "raw": "http://localhost:8081/register",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8081",
          "path": ["register"]
        }
      }
    },
    {
      "name": "Login User",
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"email\": \"test@example.com\",\n  \"password\": \"secret123\"\n}"
        },
        "url": {
          "raw": "http://localhost:8081/login",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8081",
          "path": ["login"]
        }
      }
    },
    {
      "name": "GraphQL Me",
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" },
          { "key": "Authorization", "value": "Bearer {{jwt_token}}", "type": "text" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"query\": \"query { me { id email username } }\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/query",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["query"]
        }
      }
    }
  ],
  "variable": [
    {
      "key": "jwt_token",
      "value": ""
    }
  ]
}

📌 Cara Pakai

Simpan file ini → user-service-test.postman_collection.json.

Buka Postman → Import → Upload File → pilih file ini.

Jalankan request Register User → dapat JWT.

Copy JWT ke variable jwt_token (di tab Collection → Variables).

Jalankan request GraphQL Me → harusnya dapat data user sesuai JWT.
==============================================================
File: user_service_test.postman_collection.json
{
  "info": {
    "name": "User Service Test",
    "_postman_id": "e5c89d21-8f23-44f2-9b91-cc0fbb7f1000",
    "description": "Test flow Register → Login → GraphQL Gateway",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Register User",
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n    \"email\": \"demo@example.com\",\n    \"password\": \"secret123\",\n    \"username\": \"demouser\"\n}"
        },
        "url": {
          "raw": "http://localhost:8081/register",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8081",
          "path": ["register"]
        }
      },
      "response": []
    },
    {
      "name": "Login User",
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n    \"email\": \"demo@example.com\",\n    \"password\": \"secret123\"\n}"
        },
        "url": {
          "raw": "http://localhost:8081/login",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8081",
          "path": ["login"]
        }
      },
      "response": []
    },
    {
      "name": "GraphQL - Query Me",
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" },
          {
            "key": "Authorization",
            "value": "Bearer {{jwt_token}}",
            "type": "text"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n    \"query\": \"query { me { id email username } }\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/query",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["query"]
        }
      },
      "response": []
    }
  ],
  "variable": [
    {
      "key": "jwt_token",
      "value": ""
    }
  ]
}
=====================================
📌 File: user_service_env.postman_environment.json
{
  "id": "a1f3f9a4-12e4-4f55-9e66-62d2e6db13e1",
  "name": "User Service Env",
  "values": [
    {
      "key": "base_url_user",
      "value": "http://localhost:8081",
      "enabled": true
    },
    {
      "key": "base_url_gateway",
      "value": "http://localhost:8080",
      "enabled": true
    },
    {
      "key": "jwt_token",
      "value": "",
      "enabled": true
    }
  ],
  "_postman_variable_scope": "environment",
  "_postman_exported_at": "2025-08-16T06:00:00Z",
  "_postman_exported_using": "Postman/10.23.0"
  =======================================================