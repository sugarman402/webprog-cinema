package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.FullName == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users (email, password, full_name, is_admin) VALUES ($1, $2, $3, $4) RETURNING id, created_at"
	var userID int
	var createdAt time.Time

	err = DB.QueryRow(query, req.Email, string(hashedPassword), req.FullName, req.IsAdmin).Scan(&userID, &createdAt)
	if err != nil {
		http.Error(w, "Email already registered", http.StatusBadRequest)
		return
	}

	user := User{
		ID:        userID,
		Email:     req.Email,
		FullName:  req.FullName,
		IsAdmin:   req.IsAdmin,
		CreatedAt: createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user User
	query := "SELECT id, email, password, full_name, is_admin FROM users WHERE email = $1"

	err := DB.QueryRow(query, req.Email).Scan(&user.ID, &user.Email, &user.Password, &user.FullName, &user.IsAdmin)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(user.ID, user.IsAdmin)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
		User: User{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			IsAdmin:  user.IsAdmin,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, title, description, duration, genre, poster_url FROM movies ORDER BY title")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	movies := []Movie{}
	for rows.Next() {
		var m Movie
		rows.Scan(&m.ID, &m.Title, &m.Description, &m.Duration, &m.Genre, &m.PosterURL)
		movies = append(movies, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func GetScreenings(w http.ResponseWriter, r *http.Request) {
	movieID := r.URL.Query().Get("movie_id")
	var rows *sql.Rows
	var err error

	if movieID != "" {
		rows, err = DB.Query("SELECT id, movie_id, starts_at, available_seats, price FROM screenings WHERE movie_id = $1 ORDER BY starts_at", movieID)
	} else {
		rows, err = DB.Query("SELECT id, movie_id, starts_at, available_seats, price FROM screenings ORDER BY starts_at")
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	screenings := []Screening{}
	for rows.Next() {
		var s Screening
		rows.Scan(&s.ID, &s.MovieID, &s.StartsAt, &s.AvailableSeats, &s.Price)
		screenings = append(screenings, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(screenings)
}

func GetReservations(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r)

	var query string
	var rows *sql.Rows
	var err error

	if claims.IsAdmin {
		query = `SELECT r.id, r.screening_id, r.seats, r.created_at, m.title, s.starts_at, s.price * r.seats AS price
			FROM reservations r
			JOIN screenings s ON r.screening_id = s.id
			JOIN movies m ON s.movie_id = m.id
			ORDER BY r.created_at DESC`
		rows, err = DB.Query(query)
	} else {
		query = `SELECT r.id, r.screening_id, r.seats, r.created_at, m.title, s.starts_at, s.price * r.seats AS price
			FROM reservations r
			JOIN screenings s ON r.screening_id = s.id
			JOIN movies m ON s.movie_id = m.id
			WHERE r.user_id = $1
			ORDER BY r.created_at DESC`
		rows, err = DB.Query(query, claims.UserID)
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	reservations := []ReservationResponse{}
	for rows.Next() {
		var rsv ReservationResponse
		rows.Scan(&rsv.ID, &rsv.ScreeningID, &rsv.Seats, &rsv.CreatedAt, &rsv.MovieTitle, &rsv.StartsAt, &rsv.Price)
		reservations = append(reservations, rsv)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}

func DeleteScreening(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	screeningID := vars["id"]

	query := "DELETE FROM screenings WHERE id = $1"
	result, err := DB.Exec(query, screeningID)
	if err != nil {
		http.Error(w, "Failed to delete screening", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to delete screening", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Screening not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reservationID := vars["id"]

	query := "DELETE FROM reservations WHERE id = $1"
	result, err := DB.Exec(query, reservationID)
	if err != nil {
		http.Error(w, "Failed to delete reservation", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to delete reservation", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Reservation not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func CreateReservation(w http.ResponseWriter, r *http.Request) {
	var req CreateReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.ScreeningID == 0 || req.Seats <= 0 {
		http.Error(w, "Invalid screening or seat count", http.StatusBadRequest)
		return
	}

	claims := claimsFromContext(r)
	userID := claims.UserID

	var available int
	query := "SELECT available_seats FROM screenings WHERE id = $1 FOR UPDATE"
	tx, err := DB.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	err = tx.QueryRow(query, req.ScreeningID).Scan(&available)
	if err != nil {
		http.Error(w, "Screening not found", http.StatusBadRequest)
		return
	}

	if req.Seats > available {
		http.Error(w, "Not enough available seats", http.StatusBadRequest)
		return
	}

	insertQuery := "INSERT INTO reservations (user_id, screening_id, seats, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id, created_at"
	var reservationID int
	var createdAt time.Time
	err = tx.QueryRow(insertQuery, userID, req.ScreeningID, req.Seats).Scan(&reservationID, &createdAt)
	if err != nil {
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		return
	}

	updateQuery := "UPDATE screenings SET available_seats = available_seats - $1 WHERE id = $2"
	_, err = tx.Exec(updateQuery, req.Seats, req.ScreeningID)
	if err != nil {
		http.Error(w, "Failed to save reservation", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := ReservationResponse{
		ID:          reservationID,
		ScreeningID: req.ScreeningID,
		Seats:       req.Seats,
		CreatedAt:   createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var req CreateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Description == "" || req.Duration <= 0 || req.Genre == "" {
		http.Error(w, "All required fields must be filled", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO movies (title, description, duration, genre, poster_url) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var movieID int
	err := DB.QueryRow(query, req.Title, req.Description, req.Duration, req.Genre, req.PosterURL).Scan(&movieID)
	if err != nil {
		http.Error(w, "Failed to create movie", http.StatusInternalServerError)
		return
	}

	movie := Movie{
		ID:          movieID,
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		Genre:       req.Genre,
		PosterURL:   req.PosterURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID := vars["id"]

	query := "DELETE FROM movies WHERE id = $1"
	result, err := DB.Exec(query, movieID)
	if err != nil {
		http.Error(w, "Failed to delete movie", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to delete movie", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func CreateScreening(w http.ResponseWriter, r *http.Request) {
	var req CreateScreeningRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.MovieID == 0 || req.StartsAt == "" || req.AvailableSeats <= 0 || req.Price <= 0 {
		http.Error(w, "All required fields must be filled", http.StatusBadRequest)
		return
	}

	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		startsAt, err = time.Parse("2006-01-02T15:04", req.StartsAt)
	}
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	var exists bool
	err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM movies WHERE id = $1)", req.MovieID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Movie not found", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO screenings (movie_id, starts_at, available_seats, price) VALUES ($1, $2, $3, $4) RETURNING id"
	var screeningID int
	err = DB.QueryRow(query, req.MovieID, startsAt, req.AvailableSeats, req.Price).Scan(&screeningID)
	if err != nil {
		http.Error(w, "Failed to create screening", http.StatusInternalServerError)
		return
	}

	screening := Screening{
		ID:             screeningID,
		MovieID:        req.MovieID,
		StartsAt:       startsAt,
		AvailableSeats: req.AvailableSeats,
		Price:          req.Price,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(screening)
}

type CreateReservationRequest struct {
	ScreeningID int `json:"screening_id"`
	Seats       int `json:"seats"`
}

type CreateMovieRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Genre       string `json:"genre"`
	PosterURL   string `json:"poster_url"`
}

type CreateScreeningRequest struct {
	MovieID        int     `json:"movie_id"`
	StartsAt       string  `json:"starts_at"`
	AvailableSeats int     `json:"available_seats"`
	Price          float64 `json:"price"`
}
