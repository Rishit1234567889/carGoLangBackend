package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // The underscore import is important!
	"log"
	"os"
	"time"
)

var db *sql.DB

func InitDB() {
	connStr := fmt.Sprintf("host = %s port = %s user = %s password = %s dbname = %s sslmode = disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	fmt.Println("Connecting to PostgresSQL...")
	time.Sleep(5 * time.Second)

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening in database", err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	db = database
	fmt.Println("Successfully connected to database")
}

//const (
//	host     = "localhost" //"db" //"localhost"
//	port     = 5432
//	user     = "user"
//	password = "password"
//	dbname   = "mydb"
//)

//func InitDB() {
//	// Connection string
//	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
//		host, port, user, password, dbname)
//
//	// Open the connection
//	var err error
//	db, err = sql.Open("postgres", psqlInfo)
//	if err != nil {
//		log.Fatalf("Error opening database: %q", err)
//	}
//	//defer db.Close()
//
//	// Verify the connection
//	err = db.Ping()
//	if err != nil {
//		log.Fatalf("Error verifying connection with database: %q", err)
//	}
//
//	fmt.Println("Successfully connected to the database!")
//}

func GetDB() *sql.DB {
	return db
}

func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatalf("Error in closing the Database : %v", err)
	}
}
