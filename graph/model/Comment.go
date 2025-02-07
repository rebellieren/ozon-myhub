package model

import (
	"myhub/internal/utils"
	"time"

	"github.com/google/uuid"
)

const MaxCommentLength = 2000

type Comment struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	UserID    string     `json:"userID"`
	ParentID  string     `json:"parentId"`
	PostID    string     `json:"postId"`
	User      *User      `json:"user"`
	Replies   []*Comment `json:"replies"`
	CreatedAt time.Time  `json:"createdAt"`
}

func NewComment(content, userID, postID string) (*Comment, error) {
	if len(content) > MaxCommentLength {
		return nil, utils.NewGraphQLError("длина комментария не должна превышать 2000 символов", "VALIDATION_FAILED")
	}

	commentID := uuid.New().String()

	return &Comment{
		ID:        commentID,
		Content:   content,
		UserID:    userID,
		PostID:    postID,
		Replies:   []*Comment{},
		CreatedAt: time.Now(),
	}, nil
}
