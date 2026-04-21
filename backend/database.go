package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDatabase() error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	err = createTables()
	if err != nil {
		return err
	}

	err = seedInitialData()
	if err != nil {
		return err
	}

	return nil
}

func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		full_name VARCHAR(255) NOT NULL,
		is_admin BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS movies (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		duration INTEGER NOT NULL,
		genre VARCHAR(100) NOT NULL,
		poster_url VARCHAR(255)
	);

	CREATE TABLE IF NOT EXISTS screenings (
		id SERIAL PRIMARY KEY,
		movie_id INTEGER NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
		starts_at TIMESTAMP NOT NULL,
		available_seats INTEGER NOT NULL,
		price NUMERIC(8,2) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS reservations (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		screening_id INTEGER NOT NULL REFERENCES screenings(id) ON DELETE CASCADE,
		seats INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_screenings_movie_id ON screenings(movie_id);
	CREATE INDEX IF NOT EXISTS idx_reservations_user_id ON reservations(user_id);
	`

	_, err := DB.Exec(schema)
	return err
}

func seedInitialData() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM movies").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	_, err = DB.Exec(`
	INSERT INTO movies (title, description, duration, genre, poster_url) VALUES
	('The Dark Knight', 'Batman raises the stakes in his war on crime. With the help of Lt. Jim Gordon and District Attorney Harvey Dent, Batman sets out to dismantle the remaining criminal organizations that plague the streets. The partnership proves to be effective, but they soon find themselves prey to a reign of chaos unleashed by a rising criminal mastermind known to the terrified citizens of Gotham as the Joker.', 152, 'Action', 'https://image.tmdb.org/t/p/original/vGYJRor3pCyjbaCpJKC39MpJhIT.jpg'),
	('Interstellar', 'The adventures of a group of explorers who make use of a newly discovered wormhole to surpass the limitations on human space travel and conquer the vast distances involved in an interstellar voyage.', 169, 'Sci-Fi', 'https://www.themoviedb.org/t/p/w600_and_h900_face/yQvGrMoipbRoddT0ZR8tPoR7NfX.jpg'),
	('Pirates of the Caribbean: Dead Man''s Chest', 'Captain Jack Sparrow''s got a blood debt to pay: he owes his soul to the legendary Davy Jones, ghastly Ruler of the Ocean Depths. To escape eternal servitude aboard the Flying Dutchman, ever-crafty Jack must track down the still-beating heart of Jones. But he won''t do it alone: Will Turner and Elizabeth Swann are drawn back into another one of his perilous quests—assuming they can evade execution for aiding a pirate.', 110, 'Action', 'https://www.themoviedb.org/t/p/w600_and_h900_face/uXEqmloGyP7UXAiphJUu2v2pcuE.jpg');
	`)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
	INSERT INTO screenings (movie_id, starts_at, available_seats, price) VALUES
	((SELECT id FROM movies WHERE title = 'The Dark Knight'), '2026-04-18 18:00:00', 80, 2500.00),
	((SELECT id FROM movies WHERE title = 'The Dark Knight'), '2026-04-18 20:30:00', 80, 2500.00),
	((SELECT id FROM movies WHERE title = 'Interstellar'), '2026-04-18 17:00:00', 70, 2200.00),
	((SELECT id FROM movies WHERE title = 'Pirates of the Caribbean: Dead Man''s Chest'), '2026-04-18 19:30:00', 60, 2300.00);
	`)
	return err
}
