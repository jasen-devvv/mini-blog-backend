package models

import "time"

// Article represents a blog article in the system.
//
// Fields:
//   - ID: Unique identifier for the article.
//   - Title: Title of the article (max 255 characters, required).
//   - Content: Main content of the article (required).
//   - UserID: ID of the user who created the article.
//   - User: Associated user who wrote the article.
//   - CreatedAt: Timestamp when the article was created.
//   - UpdatedAt: Timestamp when the article was last updated.
type Article struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
