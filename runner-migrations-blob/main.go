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

const (
	createDatabaseQuery = `CREATE DATABASE blob;`

	createUserTableQuery = `
	CREATE TABLE IF NOT EXISTS "User" (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(255) UNIQUE,
		email_verified TIMESTAMP,
		image VARCHAR(255),
		password VARCHAR(255),
		username VARCHAR(50) UNIQUE,
		bio VARCHAR(500),
		avatar_icon VARCHAR(50) DEFAULT 'user',
		avatar_color VARCHAR(50) DEFAULT 'cyan',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	createInterestTableQuery = `
	CREATE TABLE IF NOT EXISTS "Interest" (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(100) UNIQUE NOT NULL,
		description VARCHAR(500),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	createBlobTableQuery = `
	CREATE TABLE IF NOT EXISTS "Blob" (
		id VARCHAR(255) PRIMARY KEY,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		user_id VARCHAR(255) NOT NULL,
		CONSTRAINT fk_user_blob FOREIGN KEY (user_id) REFERENCES "User" (id) ON DELETE CASCADE
	);`

	createCommentTableQuery = `
	CREATE TABLE IF NOT EXISTS "Comment" (
		id VARCHAR(255) PRIMARY KEY,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		user_id VARCHAR(255) NOT NULL,
		blob_id VARCHAR(255) NOT NULL,
		CONSTRAINT fk_user_comment FOREIGN KEY (user_id) REFERENCES "User" (id) ON DELETE CASCADE,
		CONSTRAINT fk_blob_comment FOREIGN KEY (blob_id) REFERENCES "Blob" (id) ON DELETE CASCADE
	);`

	createLikeTableQuery = `
	CREATE TABLE IF NOT EXISTS "Like" (
		id VARCHAR(255) PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		user_id VARCHAR(255) NOT NULL,
		blob_id VARCHAR(255) NOT NULL,
		CONSTRAINT fk_user_like FOREIGN KEY (user_id) REFERENCES "User" (id) ON DELETE CASCADE,
		CONSTRAINT fk_blob_like FOREIGN KEY (blob_id) REFERENCES "Blob" (id) ON DELETE CASCADE,
		CONSTRAINT unique_user_blob_like UNIQUE (user_id, blob_id)
	);
	`

	createBlobInterestTableQuery = `
	CREATE TABLE IF NOT EXISTS "_BlobToInterest" (
		blob_id VARCHAR(255) NOT NULL,
		interest_id VARCHAR(255) NOT NULL,
		CONSTRAINT fk_blob_interest FOREIGN KEY (blob_id) REFERENCES "Blob" (id) ON DELETE CASCADE,
		CONSTRAINT fk_interest_blob FOREIGN KEY (interest_id) REFERENCES "Interest" (id) ON DELETE CASCADE,
		PRIMARY KEY (blob_id, interest_id)
	);`

	createVerificationTokenTable = `
	CREATE TABLE "VerificationToken" (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL,
		token TEXT NOT NULL,
		expiresat TIMESTAMP NOT NULL
	);`

	createSessionTableQuery = `
	CREATE TABLE IF NOT EXISTS "Session" (
		id VARCHAR(255) PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		expires TIMESTAMP NOT NULL,
		session_token TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_user_session FOREIGN KEY (user_id) REFERENCES "User" (id) ON DELETE CASCADE
	);`

	popBlobs = `
	CREATE OR REPLACE FUNCTION pop_old_blobs()
	RETURNS void AS $$
		DELETE FROM "Blob"
		WHERE created_at < NOW() - INTERVAL '24 hours';
	$$ LANGUAGE sql;`

	createViewListBlob = `
		CREATE VIEW listBlobs AS
		SELECT
			b.id,
			b.user_id,
			b.content,
			b.created_at,
			b.updated_at,
			u.username,
			u.avatar_icon,
			u.created_at AS user_created_at,
			i.name AS interest_name,
			(SELECT COUNT(*) FROM "Like" l WHERE l.blob_id = b.id) AS likes_count,
			(SELECT COUNT(*) FROM "Comment" c WHERE c.blob_id = b.id) AS comments_count
		FROM "Blob" b
		LEFT JOIN "_BlobToInterest" bi ON bi.blob_id = b.id
		LEFT JOIN "Interest" i ON bi.interest_id = i.id
		LEFT JOIN "User" u ON b.user_id = u.id
		ORDER BY b.created_at DESC`

	createCheckProfanity = `
		CREATE OR REPLACE FUNCTION check_profanity()
		RETURNS TRIGGER AS $$
		DECLARE
			banned_words TEXT[] := ARRAY['merda', 'porra', 'caralho', 'puta', 'foda', 'bosta', 'desgraça', 'cacete'];
			word TEXT;
		BEGIN
			FOREACH word IN ARRAY banned_words LOOP
				IF NEW.content ILIKE '%' || word || '%' THEN
					RAISE EXCEPTION 'Conteúdo proibido detectado';
				END IF;
			END LOOP;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`

	createCheckBlobContentTrigger = `
		CREATE TRIGGER check_blob_content
		BEFORE INSERT OR UPDATE ON "Blob"
		FOR EACH ROW
		EXECUTE FUNCTION check_profanity();
	`

	createCheckCommentContentTrigger = `
		CREATE TRIGGER check_comment_content
		BEFORE INSERT OR UPDATE ON "Comment"
		FOR EACH ROW
		EXECUTE FUNCTION check_profanity();
	`

)

func runningQueries(db *sqlx.DB, ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlobSchemas.Create")
	defer span.Finish()

	if _, err := db.ExecContext(ctx, createDatabaseQuery); err != nil {
		log.Println("Database already exists, skipping creation:", err)
	} else {
		log.Println("Database created successfully.")
	}

	log.Println("Database created successfully (or already exists).")

	defer db.Close()

	dbConnection, err := ConnectDB("")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	queries := []string{
		createUserTableQuery,
		createInterestTableQuery,
		createBlobTableQuery,
		createCommentTableQuery,
		createLikeTableQuery,
		createBlobInterestTableQuery,
		createSessionTableQuery,
	}

	for _, query := range queries {
		if _, err := dbConnection.ExecContext(ctx, query); err != nil {
			log.Fatalf("Failed to execute query: %v", err)
		}
	}

	log.Println("All tables created successfully.")

	if _, err := dbConnection.ExecContext(ctx, popBlobs); err != nil {
		log.Fatalf("Failed to create delete_old_blobs function: %v", err)
	}
	if _, err := dbConnection.ExecContext(ctx, createCheckProfanity); err != nil {
		log.Fatalf("Failed to create createCheckProfanity function: %v", err)
	}
	if _, err := dbConnection.ExecContext(ctx, createCheckBlobContentTrigger); err != nil {
		log.Fatalf("Failed to create createCheckBlobContentTrigger trigger: %v", err)
	}
	if _, err := dbConnection.ExecContext(ctx, createCheckCommentContentTrigger); err != nil {
		log.Fatalf("Failed to create createCheckCommentContentTrigger trigger: %v", err)
	}
	if _, err := dbConnection.ExecContext(ctx, createViewListBlob); err != nil {
		log.Fatalf("Failed to create createViewListBlob view: %v", err)
	}
 
	log.Println("Database, schema, and function created successfully.")
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

	fmt.Println("Connected to " + dbname)
	return db, nil
}

func main() { 
	dbConnection, err := ConnectDB("postgres")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConnection.Close()
	ctx := context.Background()
	 
	runningQueries(dbConnection, ctx)
}