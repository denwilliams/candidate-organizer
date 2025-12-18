module github.com/candidate-organizer/backend

// +heroku install ./backend/cmd/server

go 1.24.7

require (
	github.com/go-chi/chi/v5 v5.2.3
	github.com/go-chi/cors v1.2.2
)

require (
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/oauth2 v0.34.0 // indirect
)
