package middleware

import (
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func SetSessionCookieIfAbsent(w http.ResponseWriter, r *http.Request) {
	// Check if the user already has a cookie to mark session
	if _, err := r.Cookie("user_id"); err != nil {
		// Create a new UUID for the user
		newUUID := uuid.New().String()

		// Set a new cookie with the UUID
		http.SetCookie(w, &http.Cookie{
			Name:  "user_id",
			Value: newUUID,
			Path:  "/",
			// Other attributes like Expires, Secure, HttpOnly as necessary
		})
	}
}

func Cookie_middleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SetSessionCookieIfAbsent(w, r)
		log.Println(r.Method, r.URL.Path)
		f(w, r)
	}
}
