package model

import (
	"myhub/internal/utils"

	"github.com/google/uuid"
)

const MaxPostLength = 4000

type Post struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	User            *User      `json:"user"`
	UserID          string     `json:"userID"`
	Content         string     `json:"content"`
	Comments        []*Comment `json:"comments"`
	CommentsAllowed bool       `json:"commentsAllowed"`
}

func NewPost(title, content, userID string, commentsAllowed bool) (*Post, error) {
	if len(content) > MaxPostLength {
		return nil, utils.NewGraphQLError("длина поста не должна превышать 4000 символов", "VALIDATION_FAILED")
	}

	postID := uuid.New().String()

	return &Post{
		ID:              postID,
		Title:           title,
		UserID:          userID,
		Content:         content,
		Comments:        []*Comment{},
		CommentsAllowed: commentsAllowed,
	}, nil
}
