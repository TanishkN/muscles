package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Post struct {
	ID        int
	Username  string
	Content   string
	ImageURL  string
	Timestamp time.Time
}

var DB *pgxpool.Pool

// InitDB initializes the database connection.
func InitDB() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	var err error
	DB, err = pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to the database successfully.")
	return nil
}

// CreatePost inserts a new post into the database.
func CreatePost(post Post) error {
	_, err := DB.Exec(context.Background(),
		`INSERT INTO posts (username, content, image_url, timestamp) VALUES ($1, $2, $3, $4)`,
		post.Username, post.Content, post.ImageURL, post.Timestamp,
	)
	return err
}

// GetRecentPosts retrieves recent posts from the database.
func GetRecentPosts(limit int) ([]Post, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT id, username, content, image_url, timestamp FROM posts ORDER BY timestamp DESC LIMIT $1`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Username, &post.Content, &post.ImageURL, &post.Timestamp); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
