package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	err := InitDatabase()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer DB.Close()

	router := mux.NewRouter()

	router.Use(loggingMiddleware)
	router.Use(corsMiddleware)

	router.HandleFunc("/api/register", Register).Methods(http.MethodPost)
	router.HandleFunc("/api/login", Login).Methods(http.MethodPost)
	router.HandleFunc("/api/movies", GetMovies).Methods(http.MethodGet)
	router.HandleFunc("/api/screenings", GetScreenings).Methods(http.MethodGet)

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(authMiddleware)
	protected.HandleFunc("/movies", CreateMovie).Methods(http.MethodPost)
	protected.HandleFunc("/movies/{id}", DeleteMovie).Methods(http.MethodDelete)
	protected.HandleFunc("/screenings", CreateScreening).Methods(http.MethodPost)
	protected.HandleFunc("/screenings/{id}", DeleteScreening).Methods(http.MethodDelete)
	protected.HandleFunc("/reservations", GetReservations).Methods(http.MethodGet)
	protected.HandleFunc("/reservations", CreateReservation).Methods(http.MethodPost)
	protected.HandleFunc("/reservations/{id}", DeleteReservation).Methods(http.MethodDelete)

	router.HandleFunc("/api/{rest:.*}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("%s %s %d %s", r.Method, r.RequestURI, rw.status, time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control, Pragma")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
