listing-service/
├── go.mod                  # Module Go
├── go.sum                  # Dependencies Go
├── main.go                 # Entry point REST API server
├── .env.example            # Contoh environment variables (DB, PORT)
├── config/
│   └── config.go           # Load config/env (DB URL, PORT)
├── handler/
│   └── listing_handler.go  # Endpoint createListing, getListing, updateListing, placeBid
├── service/
│   ├── listing_service.go  # Logic bisnis: create listing, bid, media, comments
│   └── category_service.go # Logic kategori dan subkategori
├── repository/
│   ├── listing_repository.go   # Query DB untuk listings, bids, media, comments
│   └── category_repository.go  # Query DB untuk categories dan subcategories
├── model/
│   ├── listing.go           # Struct Listing
│   ├── category.go          # Struct Category / Subcategory
│   ├── bid.go               # Struct Bid
│   ├── media.go             # Struct Media (images/videos)
│   └── comment.go           # Struct Comment
├── middleware/
│   └── jwt.go               # JWT verification middleware untuk REST API
├── migrations/
│   └── 001_create_listings.sql # SQL script untuk tabel listings, categories, bids, media, comments
└── Dockerfile               # Dockerfile untuk build image listing-service
=====================================================
🔹 Penjelasan fungsi tiap folder/file

main.go → start REST API server, register routes
.env.example → variabel environment seperti DB URL, PORT
config/ → load dan expose konfigurasi service
handler/ → menerima request HTTP, panggil service layer, return response JSON
service/ → logic bisnis inti: create listing, place bid, manage media, comments
repository/ → akses database PostgreSQL untuk semua tabel listing-related
model/ → struct untuk mapping data listing, bid, media, comment, category
middleware/jwt.go → verifikasi JWT tiap request yang butuh auth
migrations/ → SQL script untuk setup tabel listing-service
Dockerfile → build image listing-service