// models/post.go
package models

import "time"

type Post struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Created_at time.Time `json:created_at`
}
