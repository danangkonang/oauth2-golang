package migration

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Migration struct{}

type DB struct {
	Db *sql.DB
}

var (
	Green = "\033[32m"
	Reset = "\033[0m"
)

func Connection() *DB {
	var connection string
	postgres := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	tz := "Asia%2FJakarta"
	mysql := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		tz,
	)
	switch os.Getenv("DB_DRIVER") {
	case "postgres":
		connection = postgres
	case "mysql":
		connection = mysql
	default:
		connection = postgres
	}
	db, err := sql.Open(os.Getenv("DB_DRIVER"), connection)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(5 * time.Minute)
	return &DB{Db: db}
}
