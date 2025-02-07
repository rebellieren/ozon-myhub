package storage

import (
	"fmt"
	"myhub/db"
	"myhub/graph/model"
	"myhub/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresDataStore struct {
	DB *gorm.DB
}

func NewPostgresDataStore(db *gorm.DB) *PostgresDataStore {
	return &PostgresDataStore{DB: db}
}
func (s *PostgresDataStore) GetPosts(limit, offset int32) (*model.PostPage, error) {
	var dbPosts []*db.Post
	var totalCount int64

	if err := s.DB.Model(&db.Post{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	if err := s.DB.Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&dbPosts).Error; err != nil {
		return nil, err
	}

	posts := db.ConvertPostsFromDB(dbPosts)
	return &model.PostPage{Posts: posts, TotalCount: int32(totalCount)}, nil
}

func (s *PostgresDataStore) CreateCommentForPost(userID, postID string, content string) (*model.Comment, error) {
	if err := db.ExistsUser(s.DB, userID); err != nil {
		return nil, err
	}

	if err := db.ExistsPost(s.DB, postID); err != nil {
		return nil, err
	}
	if err := db.CanCommentOnPost(s.DB, postID); err != nil {
		return nil, err
	}
	newComment, err := model.NewComment(content, userID, postID)
	if err != nil {
		return nil, utils.NewGraphQLError(err.Error(), "INVALID_POST_DATA")
	}

	dbComment := db.Comment{
		ID:      uuid.MustParse(newComment.ID),
		Content: content,
		PostID:  uuid.MustParse(newComment.PostID),
		UserID:  uuid.MustParse(newComment.UserID),
	}
	if err := db.CanCommentOnPost(s.DB, postID); err != nil {
		return nil, err
	}
	if err := s.DB.Create(&dbComment).Error; err != nil {
		return nil, utils.NewGraphQLError("Ошибка при создании комментария", "COMMENT_CREATION_FAILED")
	}

	return &model.Comment{
		ID:      dbComment.ID.String(),
		Content: dbComment.Content,
		PostID:  dbComment.PostID.String(),
		UserID:  dbComment.UserID.String(),
	}, nil
}

func (s *PostgresDataStore) CreatePost(userID, title, content string, commentsAllowed bool) (*model.Post, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, utils.NewGraphQLError(fmt.Sprintf("некорректный userID: %s", userID), "INVALID_USER_ID")
	}

	var dbUser db.User
	if err := s.DB.First(&dbUser, "id = ?", userUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewGraphQLError(fmt.Sprintf("пользователь с ID %s не найден", userID), "USER_NOT_FOUND")
		}
		return nil, utils.NewGraphQLError("не удалось найти пользователя", "USER_SEARCH_FAILED")
	}
	newPost, err := model.NewPost(title, content, userID, commentsAllowed)
	if err != nil {
		return nil, utils.NewGraphQLError(err.Error(), "INVALID_POST_DATA")
	}

	dbPost := &db.Post{
		ID:              uuid.MustParse(newPost.ID),
		UserID:          uuid.MustParse(newPost.UserID),
		Title:           newPost.Title,
		Content:         newPost.Content,
		CommentsAllowed: newPost.CommentsAllowed,
	}

	if err := s.DB.Create(dbPost).Error; err != nil {
		return nil, utils.NewGraphQLError("ошибка при создании поста", "POST_CREATION_FAILED")
	}

	return &model.Post{
		ID:              dbPost.ID.String(),
		UserID:          dbPost.UserID.String(),
		Title:           dbPost.Title,
		Content:         dbPost.Content,
		CommentsAllowed: dbPost.CommentsAllowed,
	}, nil
}

func (s *PostgresDataStore) GetUserByID(userID string) (*model.User, error) {
	var user model.User
	if err := s.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("пользователь с ID %s не найден", userID)
	}
	return &user, nil
}

func (s *PostgresDataStore) GetPost(postID string) (*model.Post, error) {
	var dbPost db.Post
	if err := s.DB.First(&dbPost, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("пост с ID %s не найден", postID)
		}
		return nil, utils.NewGraphQLError("не удалось найти пользователя", "USER_SEARCH_FAILED")
	}
	if err := s.DB.Preload("Comments").First(&dbPost, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewGraphQLError("Пост не найден", "POST_NOT_FOUND")
		}
		return nil, utils.NewGraphQLError("Ошибка при загрузке поста", "SERVER_ERROR")
	}

	return &model.Post{
		ID:              dbPost.ID.String(),
		Title:           dbPost.Title,
		Content:         dbPost.Content,
		CommentsAllowed: dbPost.CommentsAllowed,
	}, nil
}
func (s *PostgresDataStore) ToggleCommentsForPost(postID string, userID string) (*model.Post, error) {
	var dbPost db.Post

	if err := s.DB.Where("id = ? AND user_id = ?", postID, userID).First(&dbPost).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewGraphQLError("Пост не найден или у вас нет прав на изменение", "POST_NOT_FOUND")
		}
		return nil, utils.NewGraphQLError("Ошибка при поиске поста", "DATABASE_ERROR")
	}

	dbPost.CommentsAllowed = !dbPost.CommentsAllowed

	if err := s.DB.Save(&dbPost).Error; err != nil {
		return nil, utils.NewGraphQLError("Ошибка при обновлении поста", "DATABASE_ERROR")
	}

	return &model.Post{
		ID:              dbPost.ID.String(),
		Title:           dbPost.Title,
		Content:         dbPost.Content,
		CommentsAllowed: dbPost.CommentsAllowed,
	}, nil
}

func (s *PostgresDataStore) GetCommentsByPostID(postID string, limit, offset int32) ([]*model.Comment, error) {
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %w", err)
	}

	var dbComments []db.Comment
	query := `
		SELECT c.id, c.post_id, c.content, c.created_at, c.user_id
		FROM comments c
		LEFT JOIN replies r ON c.id = r.child_id
		WHERE c.post_id = ? AND r.parent_id IS NULL
		ORDER BY c.created_at DESC
		LIMIT ? OFFSET ?`

	if err := s.DB.Raw(query, postUUID, limit, offset).Scan(&dbComments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}

	comments := make([]*model.Comment, len(dbComments))
	for i, dbComment := range dbComments {
		comments[i] = &model.Comment{
			ID:      dbComment.ID.String(),
			Content: dbComment.Content,
			PostID:  dbComment.PostID.String(),
			UserID:  dbComment.UserID.String(),
		}
	}

	return comments, nil
}
func (s *PostgresDataStore) GetRepliesByCommentID(commentID string, limit int32, offset int32) ([]*model.Comment, error) {
	commentUUID, err := uuid.Parse(commentID)
	if err != nil {
		return nil, err
	}

	var dbReplies []*db.Comment

	query := `SELECT c.id, c.post_id, c.content, c.created_at, c.user_id 
              FROM comments c
              JOIN replies r ON c.id = r.child_id
              WHERE r.parent_id = ?
              ORDER BY c.created_at ASC
              LIMIT ? OFFSET ?`

	if err := s.DB.Raw(query, commentUUID, limit, offset).Scan(&dbReplies).Error; err != nil {
		return nil, err
	}

	var replies []*model.Comment
	for _, dbComment := range dbReplies {

		replies = append(replies, &model.Comment{
			ID:      dbComment.ID.String(),
			Content: dbComment.Content,
			PostID:  dbComment.PostID.String(),
			UserID:  dbComment.UserID.String(),
		})
	}

	return replies, nil
}

func (s *PostgresDataStore) CreateReplyForComment(userID, parentCommentID, content string) (*model.Comment, error) {
	if err := db.ExistsUser(s.DB, userID); err != nil {
		return nil, err
	}

	if err := db.ExistsComment(s.DB, parentCommentID); err != nil {
		return nil, err
	}

	var parentComment db.Comment
	if err := s.DB.Preload("User").First(&parentComment, "id = ?", parentCommentID).Error; err != nil {
		return nil, utils.NewGraphQLError("Родительский комментарий не найден", "PARENT_COMMENT_NOT_FOUND")
	}
	if err := db.CanCommentOnPost(s.DB, parentComment.PostID.String()); err != nil {
		return nil, utils.NewGraphQLError("Вы не можете оставлять комментарии к этому посту", "ACCESS_DENIED")
	}

	newComment := db.Comment{
		ID:      uuid.New(),
		Content: content,
		PostID:  parentComment.PostID,
		UserID:  uuid.MustParse(userID),
	}

	if err := s.DB.Create(&newComment).Error; err != nil {
		return nil, utils.NewGraphQLError("Ошибка при создании ответа", "REPLY_CREATION_FAILED")
	}

	reply := db.Reply{
		ParentID: uuid.MustParse(parentCommentID),
		ChildID:  newComment.ID,
	}
	if err := s.DB.Create(&reply).Error; err != nil {
		return nil, utils.NewGraphQLError("Ошибка при добавлении связи комментария", "REPLY_CREATION_FAILED")
	}

	if err := s.DB.Preload("User").First(&newComment, "id = ?", newComment.ID).Error; err != nil {
		return nil, utils.NewGraphQLError("Ошибка при загрузке нового комментария", "COMMENT_FETCH_FAILED")
	}

	var userModel *model.User
	if newComment.User != nil {
		userModel = &model.User{
			ID:       newComment.User.ID.String(),
			Nickname: newComment.User.Nickname,
		}
	}

	return &model.Comment{
		ID:      newComment.ID.String(),
		Content: newComment.Content,
		PostID:  newComment.PostID.String(),
		UserID:  newComment.UserID.String(),
		User:    userModel,
	}, nil
}
