package storage

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	store := NewInMemoryDataStore()
	store.CreateTestUsersIfNotExist()

	userID := "1"
	title := "Test Post"
	content := "This is a test post"
	commentsAllowed := true

	post, err := store.CreatePost(userID, title, content, commentsAllowed)
	require.NoError(t, err)
	require.NotNil(t, post)

	assert.Equal(t, userID, post.UserID)
	assert.Equal(t, title, post.Title)
	assert.Equal(t, content, post.Content)
	assert.Equal(t, commentsAllowed, post.CommentsAllowed)
}

func TestCreateCommentForPost(t *testing.T) {
	store := NewInMemoryDataStore()
	store.CreateTestUsersIfNotExist()

	userID := "1"
	post, err := store.CreatePost(userID, "Test", "Test Content", true)
	require.NoError(t, err)

	comment, err := store.CreateCommentForPost(userID, post.ID, "Test Comment")
	require.NoError(t, err)
	require.NotNil(t, comment)

	assert.Equal(t, userID, comment.UserID)
	assert.Equal(t, post.ID, comment.PostID)
	assert.Equal(t, "Test Comment", comment.Content)
}

func TestCreateReplyForComment(t *testing.T) {
	store := NewInMemoryDataStore()
	store.CreateTestUsersIfNotExist()

	userID := "1"
	post, err := store.CreatePost(userID, "Test Post", "Test Content", true)
	require.NoError(t, err)

	comment, err := store.CreateCommentForPost(userID, post.ID, "Test Comment")
	require.NoError(t, err)

	reply, err := store.CreateReplyForComment(userID, comment.ID, "Test Reply")
	require.NoError(t, err)
	require.NotNil(t, reply)

	assert.Equal(t, userID, reply.UserID)
	assert.Equal(t, comment.ID, reply.ParentID)
	assert.Equal(t, "Test Reply", reply.Content)
}

func TestGetCommentsByPostID(t *testing.T) {
	store := NewInMemoryDataStore()
	store.CreateTestUsersIfNotExist()

	userID := "1"
	log.Printf("✅ Создаю пост: %s", "Comment 1")
	post, err := store.CreatePost(userID, "Test", "Test Content", true)
	require.NoError(t, err)
	log.Printf("✅ Создаю комментарий: %s", "Comment 2")
	_, err = store.CreateCommentForPost(userID, post.ID, "Comment 1")
	require.NoError(t, err)

	_, err = store.CreateCommentForPost(userID, post.ID, "Comment 2")
	require.NoError(t, err)
	log.Printf("📝 Получаю комментарии...")
	time.Sleep(1 * time.Millisecond)
	comments, err := store.GetCommentsByPostID(post.ID, 10, 0)
	require.NoError(t, err)
	log.Printf("📌 Найдено комментариев: %d", len(comments))
	assert.Len(t, comments, 2)
}

func TestCommentingOnDisabledPost(t *testing.T) {
	store := NewInMemoryDataStore()
	store.CreateTestUsersIfNotExist()

	userID := "1"
	post, _ := store.CreatePost(userID, "Test Post", "Content", false)

	comment, err := store.CreateCommentForPost(userID, post.ID, "Test Comment")

	require.Error(t, err)
	assert.Nil(t, comment)
}
func TestToggleCommentsByNonOwner(t *testing.T) {
	store := NewInMemoryDataStore()
	store.CreateTestUsersIfNotExist()

	ownerID := "1"
	anotherUserID := "2" // ❌ Не владелец поста
	post, _ := store.CreatePost(ownerID, "Test Post", "Content", true)

	updatedPost, err := store.ToggleCommentsForPost(post.ID, anotherUserID)

	require.Error(t, err)
	assert.Nil(t, updatedPost)
}

func TestCreateCommentForNonExistentPost(t *testing.T) {
	store := NewInMemoryDataStore()
	store.CreateTestUsersIfNotExist()

	userID := "1"
	fakePostID := "nonexistent-post"

	comment, err := store.CreateCommentForPost(userID, fakePostID, "Test Comment")
	require.Error(t, err)
	assert.Nil(t, comment)
}

func TestToggleCommentsForPost(t *testing.T) {
	store := NewInMemoryDataStore()
	store.CreateTestUsersIfNotExist()

	userID := "1"
	post, err := store.CreatePost(userID, "Test", "Test Content", true)
	require.NoError(t, err)

	updatedPost, err := store.ToggleCommentsForPost(post.ID, userID)
	require.NoError(t, err)
	assert.False(t, updatedPost.CommentsAllowed)

	updatedPost, err = store.ToggleCommentsForPost(post.ID, userID)
	require.NoError(t, err)
	assert.True(t, updatedPost.CommentsAllowed)
}
func TestGetCommentsPagination(t *testing.T) {
	store := NewInMemoryDataStore()
	store.CreateTestUsersIfNotExist()

	userID := "1"
	post, _ := store.CreatePost(userID, "Test Post", "Content", true)

	store.CreateCommentForPost(userID, post.ID, "Comment 1")
	store.CreateCommentForPost(userID, post.ID, "Comment 2")
	store.CreateCommentForPost(userID, post.ID, "Comment 3")

	comments, err := store.GetCommentsByPostID(post.ID, 2, 1)

	require.NoError(t, err)
	assert.Len(t, comments, 2)
}
func TestGetNonExistentPost(t *testing.T) {
	store := NewInMemoryDataStore()
	post, err := store.GetPost("nonexistent")

	require.Error(t, err)
	assert.Nil(t, post)
}

func TestGetPost_NotFound(t *testing.T) {
	store := NewInMemoryDataStore()

	post, err := store.GetPost("nonexistent")
	assert.Nil(t, post)
	assert.Error(t, err)
}

func TestGetUserByID(t *testing.T) {
	store := NewInMemoryDataStore()
	store.CreateTestUsersIfNotExist()

	user, err := store.GetUserByID("1")
	require.NoError(t, err)
	require.NotNil(t, user)

	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "Alice", user.Nickname)
}
