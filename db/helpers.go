package db

import (
	"myhub/internal/utils"

	"gorm.io/gorm"
)

func ExistsUser(db *gorm.DB, userID string) error {
	var count int64
	var user User
	if err := db.Model(&user).Where("id = ?", userID).Count(&count).Error; err != nil {
		return utils.NewGraphQLError("Ошибка при проверке пользователя", "USER_CHECK_FAILED")
	}
	if count == 0 {
		return utils.NewGraphQLError("Пользователь с таким ID не найден", "USER_NOT_FOUND")
	}
	return nil
}

func ExistsComment(db *gorm.DB, commentID string) error {
	var count int64
	var comment Comment
	if err := db.Model(&comment).Where("id = ?", commentID).Count(&count).Error; err != nil {
		return utils.NewGraphQLError("Ошибка при проверке комментария", "COMMENT_CHECK_FAILED")
	}
	if count == 0 {
		return utils.NewGraphQLError("Комментарий с таким ID не найден", "COMMENT_NOT_FOUND")
	}
	return nil
}
func CanCommentOnPost(db *gorm.DB, postID string) error {
	var post Post
	if err := db.Select("id, comments_allowed").Where("id = ?", postID).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.NewGraphQLError("Пост с таким ID не найден", "POST_NOT_FOUND")
		}
		return utils.NewGraphQLError("Ошибка при проверке поста", "POST_CHECK_FAILED")
	}

	if !post.CommentsAllowed {
		return utils.NewGraphQLError("Комментарии к этому посту запрещены", "COMMENTS_DISABLED")
	}

	return nil
}
func ExistsPost(db *gorm.DB, postID string) error {
	var count int64
	var post Post
	if err := db.Model(&post).Where("id = ?", postID).Count(&count).Error; err != nil {
		return utils.NewGraphQLError("Ошибка при проверке поста", "POST_CHECK_FAILED")
	}
	if count == 0 {
		return utils.NewGraphQLError("Пост с таким ID не найден", "POST_NOT_FOUND")
	}
	return nil
}
