package converter

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
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

func EntityLikeToModelLikeResponse(ctx context.Context, like *entity.Like, res *model.LikeResponse) {
	res.ID = like.ID
	res.UserID = like.UserID
	res.ImageID = like.ImageID
	res.CreatedAt = like.CreatedAt
	res.UpdatedAt = like.UpdatedAt
	res.DeletedAt = like.DeletedAt
}

func EntityLikeListToModelLikeResponseList(ctx context.Context, likeList *entity.LikeList, res *model.LikeResponseList) {
	for _, like := range *likeList {
		r := model.LikeResponse{}
		EntityLikeToModelLikeResponse(ctx, &like, &r)
		*res = append(*res, r)
	}
}

func EntityCommentToModelCommentResponse(ctx context.Context, comment *entity.Comment, res *model.CommentResponse) {
	res.ID = comment.ID
	res.UserID = comment.UserID
	res.ImageID = comment.ImageID
	res.Comment = comment.Comment
	res.CreatedAt = comment.CreatedAt
	res.UpdatedAt = comment.UpdatedAt
	res.DeletedAt = comment.DeletedAt
}

func EntityCommentListToModelCommentResponseList(ctx context.Context, commentList *entity.CommentList, res *model.CommentResponseList) {
	for _, comment := range *commentList {
		r := model.CommentResponse{}
		EntityCommentToModelCommentResponse(ctx, &comment, &r)
		*res = append(*res, r)
	}
}

func ModelImageUploadedEventToModelNotifyFollowerOnUploadRequest(ctx context.Context, event *model.ImageUploadedEvent, req *model.NotifyFollowerOnUploadRequest) {
	req.UserID = event.UserID
	req.URL = event.URL
}

func ModelImageCommentedEventToModelNotifyUserImageCommentedRequest(ctx context.Context, event *model.ImageCommentedEvent, req *model.NotifyUserImageCommentedRequest) {
	req.ImageID = event.ImageID
	req.CommenterUserID = event.UserID
}

func SaramaConsumerMessageListToModelBatchUpdateImageCommentCountRequest(ctx context.Context, messages []*sarama.ConsumerMessage, req *model.BatchUpdateImageCommentCountRequest) {
	mapCounter := make(map[int64]int)
	for _, message := range messages {
		event := new(model.ImageCommentedEvent)
		if err := json.Unmarshal(message.Value, event); err != nil {
			x.Logger.WithContext(ctx).WithError(err).Warn("Failed to unmarshal image commented event")
			continue
		}
		mapCounter[event.ImageID]++
	}

	for imageID, count := range mapCounter {
		object := model.ImageIncreaseCommentCount{
			ImageID: imageID,
			Count:   count,
		}
		req.ImageIncreaseCommentCountList = append(req.ImageIncreaseCommentCountList, object)
	}
}

func ModelImageLikedEventToModelNotifyUserImageLikedRequest(ctx context.Context, event *model.ImageLikedEvent, req *model.NotifyUserImageLikedRequest) {
	req.ImageID = event.ImageID
	req.LikerUserID = event.UserID
}

func SaramaConsumerMessageListToModelBatchUpdateImageLikeCountRequest(ctx context.Context, messages []*sarama.ConsumerMessage, req *model.BatchUpdateImageLikeCountRequest) {
	mapCounter := make(map[int64]int)
	for _, message := range messages {
		event := new(model.ImageLikedEvent)
		if err := json.Unmarshal(message.Value, event); err != nil {
			x.Logger.WithContext(ctx).WithError(err).Warn("Failed to unmarshal image liked event")
			continue
		}
		mapCounter[event.ImageID]++
	}

	for imageID, count := range mapCounter {
		object := model.ImageIncreaseLikeCount{
			ImageID: imageID,
			Count:   count,
		}
		req.ImageIncreaseLikeCountList = append(req.ImageIncreaseLikeCountList, object)
	}
}
