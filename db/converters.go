package db

import "myhub/graph/model"

func ConvertPostFromDB(dbPost *Post) *model.Post {
	return &model.Post{
		ID:              dbPost.ID.String(),
		Title:           dbPost.Title,
		Content:         dbPost.Content,
		UserID:          dbPost.UserID.String(),
		CommentsAllowed: dbPost.CommentsAllowed,
	}
}

func ConvertUserFromDB(dbUser *User) *model.User {
	return &model.User{
		ID:       dbUser.ID.String(),
		Nickname: dbUser.Nickname,
	}
}

func ConvertPostsFromDB(dbPosts []*Post) []*model.Post {
	if dbPosts == nil {
		return nil
	}
	var posts []*model.Post
	for _, dbPost := range dbPosts {
		posts = append(posts, ConvertPostFromDB(dbPost))
	}
	return posts
}

func ConvertCommentFromDB(dbComment *Comment, replies []*Comment) *model.Comment {
	var convertedReplies []*model.Comment
	for _, reply := range replies {
		convertedReplies = append(convertedReplies, ConvertCommentFromDB(reply, nil))
	}

	return &model.Comment{
		ID:      dbComment.ID.String(),
		Content: dbComment.Content,
		PostID:  dbComment.PostID.String(),
		UserID:  dbComment.UserID.String(),
		Replies: convertedReplies,
	}
}
