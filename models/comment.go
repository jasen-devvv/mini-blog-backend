package models

import "time"

// Comment represents a user comment on an article.
//
// Fields:
//   - ID: Unique identifier for the comment.
//   - Content: The actual comment text (required).
//   - UserID: ID of the user who posted the comment.
//   - User: Associated user who made the comment.
//   - ArticleID: ID of the article the comment belongs to.
//   - CreatedAt: Timestamp when the comment was created.
//   - UpdatedAt: Timestamp when the comment was last updated.
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	ArticleID uint      `json:"article_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
