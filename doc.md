user-service/
├── go.mod / go.sum
├── main.go                # Entry point REST API
├── config/
│   └── config.go          # Config DB, JWT secret
├── handler/
│   ├── auth_handler.go    # Register, Login, GoogleLogin, Me
│   └── user_handler.go    # CRUD user, list, update, delete
├── service/
│   └── user_service.go    # Logic register, login, hash password, OTP
├── repository/
│   └── user_repository.go # Query DB, simpan user, ambil user by email/id
├── model/
│   └── user.go            # Struct User, OTP, session
├── middleware/
│   └── jwt.go             # JWT verification middleware untuk REST
├── migrations/
│   └── 001_create_users.sql # Schema PostgreSQL
└── Dockerfile
======================================================
🔹 Fungsi File

main.go → start REST API + router
handler/ → endpoint HTTP (register, login, me)
service/ → logic bisnis (hash password, OTP, generate JWT)
repository/ → interface ke PostgreSQL
middleware/jwt.go → verify JWT tiap request yang butuh auth
migrations/ → SQL schema PostgreSQL

🔹 Alur Request Register/Login/Me
Client → GraphQL Gateway mutation register atau login
Gateway forward ke user-service REST API
user-service:
Register → hash password, simpan user, kirim OTP
Login → verifikasi password, generate JWT
Me → verifikasi JWT, return data user
Gateway → kembalikan response ke client