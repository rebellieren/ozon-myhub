package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"context"
	"fmt"
	"myhub/graph/model"
)

// User is the resolver for the user field.
func (r *commentResolver) User(ctx context.Context, obj *model.Comment) (*model.User, error) {
	return r.Storage.GetUserByID(obj.UserID)
}

// Replies is the resolver for the replies field.
func (r *commentResolver) Replies(ctx context.Context, obj *model.Comment, limit *int32, offset *int32) ([]*model.Comment, error) {
	var limitVal int32 = 10
	var offsetVal int32 = 0

	if limit != nil {
		limitVal = *limit
	}
	if offset != nil {
		offsetVal = *offset
	}

	return r.Storage.GetRepliesByCommentID(obj.ID, limitVal, offsetVal)
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, userID string, title string, content string, commentsAllowed bool) (*model.Post, error) {
	return r.Storage.CreatePost(userID, title, content, commentsAllowed)
}

// CreateCommentForPost is the resolver for the createCommentForPost field.
func (r *mutationResolver) CreateCommentForPost(ctx context.Context, userID string, postID string, content string) (*model.Comment, error) {
	return r.Storage.CreateCommentForPost(userID, postID, content)
}

// CreateReplyForComment is the resolver for the createReplyForComment field.
func (r *mutationResolver) CreateReplyForComment(ctx context.Context, userID string, parentCommentID string, content string) (*model.Comment, error) {
	return r.Storage.CreateReplyForComment(userID, parentCommentID, content)
}

// ToggleCommentsForPost is the resolver for the toggleCommentsForPost field.
func (r *mutationResolver) ToggleCommentsForPost(ctx context.Context, postID string, userID string) (*model.Post, error) {
	return r.Storage.ToggleCommentsForPost(postID, userID)
}

// Comments is the resolver for the comments field.
func (r *postResolver) Comments(ctx context.Context, obj *model.Post, limit *int32, offset *int32) ([]*model.Comment, error) {
	return r.Storage.GetCommentsByPostID(obj.ID, *limit, *offset)
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, limit int32, offset int32) (*model.PostPage, error) {
	return r.Storage.GetPosts(limit, offset)
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	return r.Storage.GetPost(id)
}

// Comments is the resolver for the comments field in Query.
func (r *queryResolver) Comments(ctx context.Context, limit int32, offset int32, postID string) ([]*model.Comment, error) {
	return r.Storage.GetCommentsByPostID(postID, limit, offset)
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return r.Storage.GetUserByID(id)
}

// NewComment is the resolver for the newComment field.
func (r *subscriptionResolver) NewComment(ctx context.Context, postID string, userID string) (<-chan *model.Comment, error) {
	panic(fmt.Errorf("not implemented: NewComment - newComment"))
}

// Comment returns CommentResolver implementation.
func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Post returns PostResolver implementation.
func (r *Resolver) Post() PostResolver { return &postResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
