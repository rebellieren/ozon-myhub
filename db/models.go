package db

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID              uuid.UUID `gorm:"primaryKey"`
	Title           string    `gorm:"size:255;not null"`
	Content         string    `gorm:"size:4000;not null"`
	UserID          uuid.UUID `gorm:"not null;index"`
	CommentsAllowed bool      `gorm:"default:true"`
	Comments        []Comment `gorm:"foreignKey:PostID"`
	CreatedAt       time.Time `gorm:"autoCreateTime;index"`
}

type Comment struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Content   string    `gorm:"size:2000;not null"`
	PostID    uuid.UUID `gorm:"not null;index"`
	UserID    uuid.UUID `gorm:"not null;index"`
	User      *User     `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
}

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Nickname  string    `gorm:"size:50;not null;unique;index"`
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
}

type Reply struct {
	ParentID  uuid.UUID `gorm:"not null;index"`
	ChildID   uuid.UUID `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
}
