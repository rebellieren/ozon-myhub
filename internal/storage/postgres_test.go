package storage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestCreateCommentForPost_postgres(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	store := &PostgresDataStore{DB: gormDB}

	userID := uuid.New().String()
	postID := uuid.New().String()
	content := "Test Comment"

	mock.ExpectQuery(`SELECT count\(\*\) FROM "users" WHERE id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(`SELECT id, comments_allowed FROM "posts" WHERE id = \$1 ORDER BY "posts"."id" LIMIT \$2`).
		WithArgs(postID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "comments_allowed"}).AddRow(postID, true))

	mock.ExpectQuery(`SELECT comments_allowed FROM "posts" WHERE id = \$1`).
		WithArgs(postID).
		WillReturnRows(sqlmock.NewRows([]string{"comments_allowed"}).AddRow(true))

	mock.ExpectExec(`INSERT INTO "comments"`).
		WithArgs(sqlmock.AnyArg(), content, postID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	comment, err := store.CreateCommentForPost(userID, postID, content)

	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, content, comment.Content)
	assert.Equal(t, postID, comment.PostID)
	assert.Equal(t, userID, comment.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
