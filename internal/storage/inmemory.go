package storage

import (
	"log"
	"myhub/graph/model"
	"myhub/internal/utils"
	"myhub/internal/validation"
	"time"

	"github.com/google/btree"
)

type InMemoryDataStore struct {
	comments *btree.BTreeG[*model.Comment]
	users    map[string]model.User
	posts    map[string]model.Post
}

func NewInMemoryDataStore() *InMemoryDataStore {
	return &InMemoryDataStore{
		posts: make(map[string]model.Post),
		comments: btree.NewG[*model.Comment](2, func(a, b *model.Comment) bool {
			if a.CreatedAt.Equal(b.CreatedAt) {
				return a.ID < b.ID //
			}
			return a.CreatedAt.Before(b.CreatedAt)
		}),
		users: make(map[string]model.User),
	}
}

func (s *InMemoryDataStore) CreateCommentForPost(userID, postID, content string) (*model.Comment, error) {
	if err := s.validateUserExists(userID); err != nil {
		return nil, err
	}
	if err := s.validatePostExists(&postID); err != nil {
		return nil, err
	}
	if err := s.checkIfCommentsAllowed(&postID); err != nil {
		return nil, err
	}

	user, exists := s.users[userID]
	if !exists {
		log.Printf("Ошибка: пользователь %s не найден в `s.users`", userID)
		return nil, utils.NewGraphQLError("Пользователь не найден", "USER_NOT_FOUND")
	}

	comment := &model.Comment{
		ID:        utils.GenerateID(),
		Content:   content,
		PostID:    postID,
		UserID:    userID,
		User:      &user,
		Replies:   []*model.Comment{},
		CreatedAt: time.Now(),
	}
	if err := validation.ValidateComment(comment); err != nil {
		return nil, err
	}
	log.Printf("Комментарий создан: %+v", comment)
	s.comments.ReplaceOrInsert(comment)

	return comment, nil
}

func (s *InMemoryDataStore) CreateReplyForComment(userID, parentCommentID, content string) (*model.Comment, error) {
	user, exists := s.users[userID]
	if !exists {
		log.Printf("Ошибка: пользователь с ID %s не найден", userID)
		return nil, utils.NewGraphQLError("Пользователь не найден", "USER_NOT_FOUND")
	}

	var parentComment *model.Comment
	s.comments.Ascend(func(c *model.Comment) bool {
		if c.ID == parentCommentID {
			parentComment = c
			return false
		}
		return true
	})

	if parentComment == nil {
		log.Printf("Ошибка: родительский комментарий с ID %s не найден", parentCommentID)
		return nil, utils.NewGraphQLError("Родительский комментарий не найден", "COMMENT_NOT_FOUND")
	}

	reply := &model.Comment{
		ID:        utils.GenerateID(),
		Content:   content,
		PostID:    parentComment.PostID,
		ParentID:  parentCommentID,
		UserID:    userID,
		User:      &user,
		Replies:   []*model.Comment{},
		CreatedAt: time.Now(),
	}

	s.comments.ReplaceOrInsert(reply)
	parentComment.Replies = append(parentComment.Replies, reply)

	log.Printf("Ответ добавлен к комментарию %s: %+v", parentCommentID, reply)
	return reply, nil
}

func (s *InMemoryDataStore) GetCommentsByPostID(postID string, limit, offset int32) ([]*model.Comment, error) {
	var comments []*model.Comment

	s.comments.Ascend(func(c *model.Comment) bool {
		if c.PostID == postID && c.ParentID == "" {
			c.Replies = s.buildCommentTree(c.ID)
			comments = append(comments, c)
		}
		return true
	})

	start := int(offset)
	end := start + int(limit)
	if start > len(comments) {
		return []*model.Comment{}, nil
	}
	if end > len(comments) {
		end = len(comments)
	}

	return comments[start:end], nil
}

func (s *InMemoryDataStore) buildCommentTree(parentID string) []*model.Comment {
	var replies []*model.Comment

	s.comments.Ascend(func(c *model.Comment) bool {
		if c.ParentID == parentID {

			c.Replies = s.buildCommentTree(c.ID)
			replies = append(replies, c)
		}
		return true
	})

	return replies
}
func (s *InMemoryDataStore) CreatePost(userID, title, content string, commentsAllowed bool) (*model.Post, error) {
	postID := utils.GenerateID()
	user, exists := s.users[userID]
	if !exists {
		return nil, utils.NewGraphQLError("Пользователь не найден", "USER_NOT_FOUND")
	}

	post := model.Post{
		ID:              postID,
		Title:           title,
		Content:         content,
		CommentsAllowed: commentsAllowed,
		UserID:          userID,
		User:            &user,
	}

	s.posts[post.ID] = post
	log.Printf(" Пост создан: %+v", post)
	return &post, nil
}
func (s *InMemoryDataStore) GetPosts(limit, offset int32) (*model.PostPage, error) {
	if s.posts == nil {
		s.posts = make(map[string]model.Post)
	}

	postList := s.getAllPosts()
	if len(postList) == 0 {
		return &model.PostPage{Posts: []*model.Post{}, TotalCount: 0, HasNextPage: false}, nil
	}

	paginatedPosts, hasNextPage := utils.Paginate(postList, limit, offset)

	return &model.PostPage{
		Posts:       toPointerSlice(paginatedPosts),
		TotalCount:  int32(len(postList)),
		HasNextPage: hasNextPage,
	}, nil
}
func (s *InMemoryDataStore) GetUserByID(userID string) (*model.User, error) {
	if err := s.validateUserExists(userID); err != nil {
		return nil, err
	}
	user := s.users[userID]

	return &user, nil
}
func (s *InMemoryDataStore) GetRepliesByCommentID(commentID string, limit int32, offset int32) ([]*model.Comment, error) {
	var replies []*model.Comment

	s.comments.Ascend(func(c *model.Comment) bool {
		if c.ParentID == commentID {
			replies = append(replies, c)
		}
		return true
	})

	start := int(offset)
	end := start + int(limit)

	if start > len(replies) {
		return []*model.Comment{}, nil
	}
	if end > len(replies) {
		end = len(replies)
	}

	return replies[start:end], nil
}

func (s *InMemoryDataStore) getAllPosts() []model.Post {
	postList := make([]model.Post, 0, len(s.posts))
	for _, post := range s.posts {
		postList = append(postList, post)
	}
	return postList
}

func toPointerSlice(posts []model.Post) []*model.Post {
	postPointers := make([]*model.Post, len(posts))
	for i := range posts {
		postPointers[i] = &posts[i]
	}
	return postPointers
}

func (s *InMemoryDataStore) GetPost(postID string) (*model.Post, error) {
	post, exists := s.posts[postID]
	if !exists {
		log.Printf("Пост с ID %s не найден", postID)
		return nil, utils.NewGraphQLError("Пост не найден", "POST_NOT_FOUND")
	}
	if post.UserID != "" {
		user, userExists := s.users[post.UserID]
		if userExists {
			post.User = &user
		} else {
			post.User = &model.User{ID: "unknown", Nickname: "Unknown User"}
		}
	}

	return &post, nil
}

func (s *InMemoryDataStore) validateUserExists(userID string) error {

	if _, exists := s.users[userID]; !exists {
		log.Printf(" Ошибка: пользователь %s не найден", userID)
		return utils.NewGraphQLError("Пользователь не найден", "USER_NOT_FOUND")
	}

	return nil
}

func (s *InMemoryDataStore) validatePostExists(postID *string) error {
	if _, exists := s.posts[*postID]; !exists {
		return utils.NewGraphQLError("Пост не найден", "POST_NOT_FOUND")
	}
	return nil
}
func (s *InMemoryDataStore) checkIfCommentsAllowed(postID *string) error {
	if err := s.validatePostExists(postID); err != nil {
		log.Printf("Ошибка: %v", err)
		return err
	}
	post := s.posts[*postID]
	if !post.CommentsAllowed {
		return utils.NewGraphQLError("Вы не можете оставлять комментарии к этому посту", "ACCESS_DENIED")
	}
	return nil
}

func (s *InMemoryDataStore) ToggleCommentsForPost(postID string, userID string) (*model.Post, error) {
	if err := s.validateUserExists(userID); err != nil {
		log.Printf("Ошибка: %v", err)
		return nil, err
	}

	if err := s.validatePostExists(&postID); err != nil {
		log.Printf("Ошибка: %v", err)
		return nil, err
	}

	post := s.posts[postID]
	if post.UserID != userID {
		return nil, utils.NewGraphQLError("Вы не можете менять статус данному посту", "ACCESS_DENIED")
	}
	post.CommentsAllowed = !post.CommentsAllowed

	s.posts[postID] = post

	log.Printf(" Статус комментариев у поста %s изменён: %v", postID, post.CommentsAllowed)
	return &post, nil
}

func (s *InMemoryDataStore) CreateTestUsersIfNotExist() {
	if len(s.users) > 0 {
		return
	}
	users := []model.User{
		{ID: "1", Nickname: "Alice"},
		{ID: "2", Nickname: "Bob"},
		{ID: "3", Nickname: "Charlie"},
	}

	for _, user := range users {
		s.users[user.ID] = user
	}

	log.Println("✅ Тестовые пользователи успешно добавлены в локальное хранилище")
	log.Println("✅ Список всех пользователей:")
	for _, user := range users {
		log.Printf("Пользователь: %s (ID: %s)\n", user.Nickname, user.ID)
	}
}
