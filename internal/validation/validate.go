package validation

import (
	"myhub/graph/model"
	"myhub/internal/utils"
)

func ValidateComment(comment *model.Comment) error {
	if comment == nil {
		return utils.NewGraphQLError("Комментарий не может быть nil", "COMMENT_IS_NULL")
	}
	if comment.ID == "" {
		return utils.NewGraphQLError("ID комментария не может быть пустым", "INVALID_COMMENT_ID")
	}
	if comment.Content == "" {
		return utils.NewGraphQLError("Содержимое комментария не может быть пустым", "EMPTY_COMMENT_CONTENT")
	}
	if len(comment.Content) > 2000 {
		return utils.NewGraphQLError("Содержимое комментария превышает 2000 символов", "COMMENT_CONTENT_TOO_LONG")
	}
	if comment.PostID == "" {
		return utils.NewGraphQLError("Комментарий должен быть привязан к посту ", "MISSING_POST_ID")
	}
	if comment.UserID == "" {
		return utils.NewGraphQLError("Комментарий должен быть привязан к пользователю )", "MISSING_USER_ID")
	}
	return nil
}
