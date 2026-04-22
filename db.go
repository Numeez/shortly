package shortly

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type DB struct {
	Conn *sql.DB
}

type Database interface {
	InsertUrl(url string) (int64, error)
	GetUrl(id int64) (string, error)
}

func connectDB() (*DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	if dbHost == "" {
		return nil, errors.New("Environment variable DB_PORT is not set")
	}
	if dbHost == "" {
		return nil, errors.New("Environment variable DB_HOST is not set")
	}
	if dbUser == "" {
		return nil, errors.New("Environment variable DB_USER is not set")
	}

	if dbPassword == "" {
		return nil, errors.New("Environment variable DB_PASSWORD is not set ")
	}
	if dbName == "" {
		return nil, errors.New("Environment variable DB_NAME is not set")
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	sqlDB, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return nil, err
	}
	log.Println("Connected to DB")
	return &DB{
		Conn: sqlDB,
	}, nil

}

func (db *DB) InsertUrl(url string) (int64, error) {
	var id int64
	expiryTime := time.Now().Add(30 * 24 * time.Hour)
	err := db.Conn.QueryRow(
		`
		INSERT INTO urls (long_url, expiry_time, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (long_url)
		DO UPDATE SET long_url = EXCLUDED.long_url
		RETURNING id;
	`,
		url, expiryTime, time.Now(),
	).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, err
}

func (db *DB) GetUrl(id int64) (string, error) {
	var result string
	var expiryTime time.Time
	err := db.Conn.QueryRow(
		` SELECT long_url, expiry_time FROM urls WHERE id = $1`, id,
	).Scan(&result, &expiryTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("url not found")
		}
		return result, err
	}
	if time.Now().After(expiryTime) {
		return "", fmt.Errorf("url has expired")
	}
	return result, nil
}
