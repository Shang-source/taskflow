package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// DB Global database connection object
var DB *sql.DB

func init() {
	// Load the .env file during package initialization (or alternatively in main.go)
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using default environment variables.")
	}
}

// ConnectDB connect database
func ConnectDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Database fail to connect: %v", err)
	}
	// Test connectivity
	if err = db.Ping(); err != nil {
		log.Fatalf("Database cannot Ping: %v", err)
	}

	log.Println("Database connect success!")
	DB = db

	// execute table creation and other operations
	initTables()
}

// initTables initialize table structure
func initTables() {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	);
	`
	_, err := DB.Exec(createUserTable)
	if err != nil {
		log.Printf("Create users tables error: %v\n", err)
	}

	createTaskTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		is_done BOOLEAN NOT NULL DEFAULT false,
	    due_date DATETIME NULL,
      reminded BOOLEAN NOT NULL DEFAULT FALSE,
		user_id BIGINT,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	_, err = DB.Exec(createTaskTable)
	if err != nil {
		log.Printf("Create tasks tables error: %v\n", err)

	}

}
