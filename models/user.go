package models

import "time"

// User represents a registered user in the system.
//
// Fields:
//   - ID: Unique identifier for the user.
//   - Username: User's unique username (max 100 characters, required).
//   - Email: User's unique email address (max 255 characters, required).
//   - Password: Hashed password of the user (hidden from JSON responses).
//   - CreatedAt: Timestamp when the user account was created.
//   - UpdatedAt: Timestamp when the user account was last updated.
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:100;not null;unique" json:"username"`
	Email     string    `gorm:"size:255;not null;unique" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"` // Hidden from JSON responses
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
