package main

import (
	"fmt"
	"os"
	"log"
	"context"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func runningQueries(db *sqlx.DB, ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobSchemas.DeleteOldBlobs")
	defer span.Finish()

	query := "SELECT pop_old_blobs();"

	if _, err := db.ExecContext(ctx, query); err != nil {
		log.Println("Failed to execute pop_old_blobs function:", err)
		return
	}
 
	log.Println("Old blobs deleted successfully!")
}

func ConnectDB(dbname string) (*sqlx.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables.")
	}

	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWD")
	)
	
	if dbname == "" {
		dbname = os.Getenv("DB_DATABASE")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Connected to: " + dbname)
	return db, nil
}

func main() { 
	dbConnection, err := ConnectDB("")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConnection.Close()
	ctx := context.Background()
	
	runningQueries(dbConnection, ctx)
}