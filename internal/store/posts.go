package store

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"log"
	"time"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostStore struct {
	db *sql.DB
}

func (p *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, tags, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := p.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		pq.Array(post.Tags),
		post.UserID,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		log.Printf("Failed to create post: %v, Error: %v", post, err)
		return err
	}

	log.Printf("Post created successfully: %+v", post)
	return nil
}
