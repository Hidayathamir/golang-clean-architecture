package converter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func ModelUploadImageRequestToModelS3UploadImageRequest(ctx context.Context, req *model.UploadImageRequest, s3UploadImgReq *model.S3UploadImageRequest) error {
	timenow := time.Now().Unix()
	userAuth := ctxuserauth.Get(ctx)
	safeFilename := strings.ReplaceAll(req.File.Filename, " ", "_")
	s3UploadImgReq.Key = fmt.Sprintf("%s/%v_%s", userAuth.Username, timenow, safeFilename)
	file, err := req.File.Open()
	if err != nil {
		return errkit.AddFuncName(err)
	}
	defer x.LogIfErr(ctx, file.Close())
	s3UploadImgReq.Body = file
	return nil
}

func EntityImageToModelImageUploadedEvent(ctx context.Context, image *entity.Image, event *model.ImageUploadedEvent) {
	event.ID = image.ID
	event.UserID = image.UserID
	event.URL = image.URL
	event.LikeCount = image.LikeCount
	event.CommentCount = image.CommentCount
	event.CreatedAt = image.CreatedAt
	event.UpdatedAt = image.UpdatedAt
	event.DeletedAt = image.DeletedAt
}

func ModelLikeImageRequestToEntityLike(ctx context.Context, req *model.LikeImageRequest, like *entity.Like) {
	userAuth := ctxuserauth.Get(ctx)
	like.UserID = userAuth.ID
	like.ImageID = req.ImageID
}

func EntityLikeToModelImageLikedEvent(ctx context.Context, like *entity.Like, event *model.ImageLikedEvent) {
	event.ID = like.ID
	event.UserID = like.UserID
	event.ImageID = like.ImageID
	event.CreatedAt = like.CreatedAt
	event.UpdatedAt = like.UpdatedAt
	event.DeletedAt = like.DeletedAt
}

func ModelCommentImageRequestToEntityComment(ctx context.Context, req *model.CommentImageRequest, comment *entity.Comment) {
	userAuth := ctxuserauth.Get(ctx)
	comment.UserID = userAuth.ID
	comment.ImageID = req.ImageID
	comment.Comment = req.Comment
}

func EntityCommentToModelImageCommentedEvent(ctx context.Context, comment *entity.Comment, event *model.ImageCommentedEvent) {
	event.ID = comment.ID
	event.UserID = comment.UserID
	event.ImageID = comment.ImageID
	event.Comment = comment.Comment
	event.CreatedAt = comment.CreatedAt
	event.UpdatedAt = comment.UpdatedAt
	event.DeletedAt = comment.DeletedAt
}

func EntityImageToModelImageResponse(ctx context.Context, image *entity.Image, res *model.ImageResponse) {
	res.ID = image.ID
	res.UserID = image.UserID
	res.URL = image.URL
	res.LikeCount = image.LikeCount
	res.CommentCount = image.CommentCount
	res.CreatedAt = image.CreatedAt
	res.UpdatedAt = image.UpdatedAt
	res.DeletedAt = image.DeletedAt
}
