package storage

import (
	"myhub/graph/model"
)

type DataStore interface {
	GetPosts(limit, offset int32) (*model.PostPage, error)
	CreatePost(userID, title, content string, commentsAllowed bool) (*model.Post, error)
	CreateCommentForPost(userID string, postID string, content string) (*model.Comment, error)
	CreateReplyForComment(userID, parentCommentID, content string) (*model.Comment, error)
	GetUserByID(userID string) (*model.User, error)
	GetCommentsByPostID(postID string, limit, offset int32) ([]*model.Comment, error)
	GetRepliesByCommentID(commentID string, limit int32, offset int32) ([]*model.Comment, error)
	GetPost(postID string) (*model.Post, error)
	ToggleCommentsForPost(postID string, userID string) (*model.Post, error)
}
