package imageusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/imageusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/pagination"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestImageUsecaseImpl_GetLike_Success_FirstPage(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	likeRepository := &mock.LikeRepositoryMock{}

	u := &imageusecase.ImageUsecaseImpl{
		DB:             gormDB,
		LikeRepository: likeRepository,
	}

	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	likeRepository.FindByImageIDFunc = func(ctx context.Context, db *gorm.DB, likeList *entity.LikeList, imageID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
		require.Equal(t, int64(100), imageID)
		require.Equal(t, int64(0), cursorID)
		require.Equal(t, dto.DefaultPageLimit+1, limit)
		*likeList = entity.LikeList{
			{ID: 30, ImageID: 100, CreatedAt: base.Add(2 * time.Hour)},
			{ID: 20, ImageID: 100, CreatedAt: base.Add(1 * time.Hour)},
		}
		return nil
	}

	res, err := u.GetLike(context.Background(), dto.GetLikeRequest{ImageID: 100})

	require.Nil(t, err)
	require.Len(t, res.Data, 2)
	require.Empty(t, res.NextCursor)
	require.Equal(t, int64(30), res.Data[0].ID)
}

func TestImageUsecaseImpl_GetLike_Success_NextPage(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	likeRepository := &mock.LikeRepositoryMock{}

	u := &imageusecase.ImageUsecaseImpl{
		DB:             gormDB,
		LikeRepository: likeRepository,
	}

	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cursor := pagination.Encode(pagination.Cursor{CreatedAt: base, ID: 100})

	likeRepository.FindByImageIDFunc = func(ctx context.Context, db *gorm.DB, likeList *entity.LikeList, imageID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
		require.Equal(t, int64(100), cursorID)
		require.True(t, base.Equal(cursorCreatedAt))
		require.Equal(t, 3, limit)
		*likeList = entity.LikeList{
			{ID: 3, ImageID: 100, CreatedAt: base.Add(-time.Hour)},
			{ID: 2, ImageID: 100, CreatedAt: base.Add(-2 * time.Hour)},
			{ID: 1, ImageID: 100, CreatedAt: base.Add(-3 * time.Hour)},
		}
		return nil
	}

	req := dto.GetLikeRequest{ImageID: 100, Page: dto.PageRequest{Limit: 2, Cursor: cursor}}
	res, err := u.GetLike(context.Background(), req)

	require.Nil(t, err)
	require.Len(t, res.Data, 2)

	next, ok, err := pagination.Decode(res.NextCursor)
	require.Nil(t, err)
	require.True(t, ok)
	require.Equal(t, int64(2), next.ID)
}

func TestImageUsecaseImpl_GetLike_Fail_InvalidCursor(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &imageusecase.ImageUsecaseImpl{DB: gormDB}

	req := dto.GetLikeRequest{ImageID: 100, Page: dto.PageRequest{Cursor: "not-a-valid-cursor"}}
	res, err := u.GetLike(context.Background(), req)

	require.Equal(t, dto.LikesResponse{}, res)
	require.NotNil(t, err)
}

func TestImageUsecaseImpl_GetComment_Success_FirstPage(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	commentRepository := &mock.CommentRepositoryMock{}

	u := &imageusecase.ImageUsecaseImpl{
		DB:                gormDB,
		CommentRepository: commentRepository,
	}

	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	commentRepository.FindByImageIDFunc = func(ctx context.Context, db *gorm.DB, commentList *entity.CommentList, imageID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
		require.Equal(t, int64(100), imageID)
		require.Equal(t, int64(0), cursorID)
		require.Equal(t, dto.DefaultPageLimit+1, limit)
		*commentList = entity.CommentList{
			{ID: 30, ImageID: 100, Comment: "c30", CreatedAt: base.Add(2 * time.Hour)},
			{ID: 20, ImageID: 100, Comment: "c20", CreatedAt: base.Add(1 * time.Hour)},
		}
		return nil
	}

	res, err := u.GetComment(context.Background(), dto.GetCommentRequest{ImageID: 100})

	require.Nil(t, err)
	require.Len(t, res.Data, 2)
	require.Empty(t, res.NextCursor)
	require.Equal(t, int64(30), res.Data[0].ID)
}

func TestImageUsecaseImpl_GetComment_Success_NextPage(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	commentRepository := &mock.CommentRepositoryMock{}

	u := &imageusecase.ImageUsecaseImpl{
		DB:                gormDB,
		CommentRepository: commentRepository,
	}

	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cursor := pagination.Encode(pagination.Cursor{CreatedAt: base, ID: 100})

	commentRepository.FindByImageIDFunc = func(ctx context.Context, db *gorm.DB, commentList *entity.CommentList, imageID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
		require.Equal(t, int64(100), cursorID)
		require.True(t, base.Equal(cursorCreatedAt))
		require.Equal(t, 3, limit)
		*commentList = entity.CommentList{
			{ID: 3, ImageID: 100, CreatedAt: base.Add(-time.Hour)},
			{ID: 2, ImageID: 100, CreatedAt: base.Add(-2 * time.Hour)},
			{ID: 1, ImageID: 100, CreatedAt: base.Add(-3 * time.Hour)},
		}
		return nil
	}

	req := dto.GetCommentRequest{ImageID: 100, Page: dto.PageRequest{Limit: 2, Cursor: cursor}}
	res, err := u.GetComment(context.Background(), req)

	require.Nil(t, err)
	require.Len(t, res.Data, 2)

	next, ok, err := pagination.Decode(res.NextCursor)
	require.Nil(t, err)
	require.True(t, ok)
	require.Equal(t, int64(2), next.ID)
}

func TestImageUsecaseImpl_GetComment_Fail_InvalidCursor(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &imageusecase.ImageUsecaseImpl{DB: gormDB}

	req := dto.GetCommentRequest{ImageID: 100, Page: dto.PageRequest{Cursor: "not-a-valid-cursor"}}
	res, err := u.GetComment(context.Background(), req)

	require.Equal(t, dto.CommentsResponse{}, res)
	require.NotNil(t, err)
}
