package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"web-crawler/scraper"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DBManager struct {
	db       *sql.DB
	host     string
	port     string
	user     string
	password string
	dbname   string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ConnectToDB() (*DBManager, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected!")
	return &DBManager{db: db, host: host, port: port, user: user, password: password, dbname: dbname}, nil
}

func (manager *DBManager) StoreInDB(items []scraper.ScrapedItem) {
	sqlStatement := `
	INSERT INTO books (title, price, availability)
	VALUES ($1, $2, $3)
	RETURNING id`
	for _, item := range items {
		_, err := manager.db.Exec(sqlStatement, item.Title, item.Price, item.Availability)
		if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
		}
	}
}
