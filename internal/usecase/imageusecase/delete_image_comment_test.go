package imageusecase_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/imageusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newDeleteImageCommentUsecase(t *testing.T) (*imageusecase.ImageUsecaseImpl, *mock.CommentRepositoryMock, *mock.ImageRepositoryMock, *mock.ImageProducerMock, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, sqlMockDB := newFakeDB(t)
	commentRepository := &mock.CommentRepositoryMock{}
	imageRepository := &mock.ImageRepositoryMock{}
	imageProducer := &mock.ImageProducerMock{}
	u := &imageusecase.ImageUsecaseImpl{
		DB:                gormDB,
		CommentRepository: commentRepository,
		ImageRepository:   imageRepository,
		ImageProducer:     imageProducer,
	}
	return u, commentRepository, imageRepository, imageProducer, sqlMockDB
}

func TestImageUsecaseImpl_DeleteImageComment_Success_CommentOwner(t *testing.T) {
	u, commentRepository, imageRepository, imageProducer, sqlMockDB := newDeleteImageCommentUsecase(t)

	commentRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, comment *entity.Comment, id int64) error {
		*comment = entity.Comment{ID: 50, UserID: 1, ImageID: 100}
		return nil
	}
	imageRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, image *entity.Image, id int64) error {
		*image = entity.Image{ID: 100, UserID: 2}
		return nil
	}
	var deleted bool
	commentRepository.DeleteByIDFunc = func(ctx context.Context, db *gorm.DB, comment *entity.Comment) error {
		deleted = true
		return nil
	}
	var sentEvent *dto.ImageCommentedEvent
	imageProducer.SendImageCommentedFunc = func(ctx context.Context, db *gorm.DB, event *dto.ImageCommentedEvent) error {
		sentEvent = event
		return nil
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 1})

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectCommit()

	err := u.DeleteImageComment(ctx, dto.DeleteImageCommentRequest{ImageID: 100, CommentID: 50})

	require.Nil(t, err)
	require.True(t, deleted)
	require.NotNil(t, sentEvent)
	require.Equal(t, -1, sentEvent.Count)
}

func TestImageUsecaseImpl_DeleteImageComment_Success_ImageOwner(t *testing.T) {
	u, commentRepository, imageRepository, imageProducer, sqlMockDB := newDeleteImageCommentUsecase(t)

	commentRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, comment *entity.Comment, id int64) error {
		*comment = entity.Comment{ID: 50, UserID: 1, ImageID: 100}
		return nil
	}
	imageRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, image *entity.Image, id int64) error {
		*image = entity.Image{ID: 100, UserID: 9}
		return nil
	}
	commentRepository.DeleteByIDFunc = func(ctx context.Context, db *gorm.DB, comment *entity.Comment) error {
		return nil
	}
	imageProducer.SendImageCommentedFunc = func(ctx context.Context, db *gorm.DB, event *dto.ImageCommentedEvent) error {
		return nil
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 9})

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectCommit()

	err := u.DeleteImageComment(ctx, dto.DeleteImageCommentRequest{ImageID: 100, CommentID: 50})

	require.Nil(t, err)
}

func TestImageUsecaseImpl_DeleteImageComment_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &imageusecase.ImageUsecaseImpl{DB: gormDB}

	err := u.DeleteImageComment(context.Background(), dto.DeleteImageCommentRequest{})

	require.NotNil(t, err)
}

func TestImageUsecaseImpl_DeleteImageComment_Fail_CommentNotOnImage(t *testing.T) {
	u, commentRepository, _, _, _ := newDeleteImageCommentUsecase(t)

	commentRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, comment *entity.Comment, id int64) error {
		*comment = entity.Comment{ID: 50, UserID: 1, ImageID: 999}
		return nil
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 1})

	err := u.DeleteImageComment(ctx, dto.DeleteImageCommentRequest{ImageID: 100, CommentID: 50})

	require.NotNil(t, err)
	require.Equal(t, 404, errkit.GetHTTPError(err).HTTPCode)
}

func TestImageUsecaseImpl_DeleteImageComment_Fail_Forbidden(t *testing.T) {
	u, commentRepository, imageRepository, _, _ := newDeleteImageCommentUsecase(t)

	commentRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, comment *entity.Comment, id int64) error {
		*comment = entity.Comment{ID: 50, UserID: 1, ImageID: 100}
		return nil
	}
	imageRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, image *entity.Image, id int64) error {
		*image = entity.Image{ID: 100, UserID: 2}
		return nil
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 7})

	err := u.DeleteImageComment(ctx, dto.DeleteImageCommentRequest{ImageID: 100, CommentID: 50})

	require.NotNil(t, err)
	require.Equal(t, 403, errkit.GetHTTPError(err).HTTPCode)
}

func TestImageUsecaseImpl_DeleteImageComment_Fail_FindComment(t *testing.T) {
	u, commentRepository, _, _, _ := newDeleteImageCommentUsecase(t)

	commentRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, comment *entity.Comment, id int64) error {
		return assert.AnError
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 1})

	err := u.DeleteImageComment(ctx, dto.DeleteImageCommentRequest{ImageID: 100, CommentID: 50})

	require.ErrorIs(t, err, assert.AnError)
}
