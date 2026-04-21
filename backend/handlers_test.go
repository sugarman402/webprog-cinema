package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()
	DB = db

	// Mock the database query
	mock.ExpectQuery(`INSERT INTO users \(email, password, full_name, is_admin\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id, created_at`).
		WithArgs("test@example.com", sqlmock.AnyArg(), "Test User", false).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))

	reqBody := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
		IsAdmin:  false,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Register(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response User
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", response.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestRegisterInvalidData(t *testing.T) {
	reqBody := RegisterRequest{
		Email:    "",
		Password: "password123",
		FullName: "Test User",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()
	DB = db

	// Generate hash for "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mock.ExpectQuery(`SELECT id, email, password, full_name, is_admin FROM users WHERE email = \$1`).
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "full_name", "is_admin"}).
			AddRow(1, "test@example.com", string(hashedPassword), "Test User", false))

	reqBody := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response AuthResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.User.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", response.User.Email)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()
	DB = db

	mock.ExpectQuery(`SELECT id, email, password, full_name, is_admin FROM users WHERE email = \$1`).
		WithArgs("wrong@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "full_name", "is_admin"}))

	reqBody := LoginRequest{
		Email:    "wrong@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestGetMovies(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()
	DB = db

	mock.ExpectQuery(`SELECT id, title, description, duration, genre, poster_url FROM movies ORDER BY title`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "duration", "genre", "poster_url"}).
			AddRow(1, "Test Movie", "Description", 120, "Action", ""))

	req := httptest.NewRequest("GET", "/api/movies", nil)
	w := httptest.NewRecorder()

	GetMovies(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var movies []Movie
	json.Unmarshal(w.Body.Bytes(), &movies)
	if len(movies) != 1 || movies[0].Title != "Test Movie" {
		t.Errorf("Expected 1 movie with title 'Test Movie', got %v", movies)
	}
}

func TestCreateMovie(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()
	DB = db

	mock.ExpectQuery(`INSERT INTO movies \(title, description, duration, genre, poster_url\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
		WithArgs("New Movie", "Description", 100, "Drama", "").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	reqBody := CreateMovieRequest{
		Title:       "New Movie",
		Description: "Description",
		Duration:    100,
		Genre:       "Drama",
		PosterURL:   "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/movies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateMovie(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var movie Movie
	json.Unmarshal(w.Body.Bytes(), &movie)
	if movie.Title != "New Movie" {
		t.Errorf("Expected title 'New Movie', got %s", movie.Title)
	}
}

func TestCreateMovieInvalidData(t *testing.T) {
	reqBody := CreateMovieRequest{
		Title:       "",
		Description: "Description",
		Duration:    100,
		Genre:       "Drama",
		PosterURL:   "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/movies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateMovie(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteMovie(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()
	DB = db

	mock.ExpectExec(`DELETE FROM movies WHERE id = \$1`).
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("DELETE", "/api/movies/1", nil)
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	DeleteMovie(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestDeleteMovieNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()
	DB = db

	mock.ExpectExec(`DELETE FROM movies WHERE id = \$1`).
		WithArgs("999").
		WillReturnResult(sqlmock.NewResult(0, 0))

	req := httptest.NewRequest("DELETE", "/api/movies/999", nil)
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"id": "999"})

	DeleteMovie(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}