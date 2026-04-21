package main

import "time"

// User model
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FullName  string    `json:"full_name"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

// Movie model
type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Genre       string `json:"genre"`
	PosterURL   string `json:"poster_url"`
}

// Screening model
type Screening struct {
	ID             int       `json:"id"`
	MovieID        int       `json:"movie_id"`
	StartsAt       time.Time `json:"starts_at"`
	AvailableSeats int       `json:"available_seats"`
	Price          float64   `json:"price"`
}

// Reservation model
type Reservation struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ScreeningID int       `json:"screening_id"`
	Seats       int       `json:"seats"`
	CreatedAt   time.Time `json:"created_at"`
}

// ReservationResponse for querying user reservations
type ReservationResponse struct {
	ID          int       `json:"id"`
	ScreeningID int       `json:"screening_id"`
	Seats       int       `json:"seats"`
	CreatedAt   time.Time `json:"created_at"`
	MovieTitle  string    `json:"movie_title"`
	StartsAt    time.Time `json:"starts_at"`
	Price       float64   `json:"price"`
}

// LoginRequest for the login endpoint
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest for the registration endpoint
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	IsAdmin  bool   `json:"is_admin"`
}

// AuthResponse for authentication responses
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
