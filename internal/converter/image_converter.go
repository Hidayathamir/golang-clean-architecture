package converter

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/twmb/franz-go/pkg/kgo"
)

func DtoUploadImageRequestToDtoS3UploadImageRequest(ctx context.Context, req dto.UploadImageRequest, s3UploadImgReq *dto.S3UploadImageRequest) error {
	timenow := time.Now().Unix()
	userAuth := ctxuserauth.Get(ctx)
	safeFilename := strings.ReplaceAll(req.File.Filename, " ", "_")
	s3UploadImgReq.Key = fmt.Sprintf("%s/%v_%s", userAuth.Username, timenow, safeFilename)
	file, err := req.File.Open()
	if err != nil {
		return errkit.AddFuncName(err)
	}
	defer x.LogIfErrForDeferContext(ctx, file.Close)
	s3UploadImgReq.Body = file
	return nil
}

func EntityImageToDtoImageUploadedEvent(image entity.Image, event *dto.ImageUploadedEvent) {
	event.ID = image.ID
	event.UserID = image.UserID
	event.Caption = image.Caption
	event.URL = image.URL
	event.LikeCount = image.LikeCount
	event.CommentCount = image.CommentCount
	event.CreatedAt = image.CreatedAt
	event.UpdatedAt = image.UpdatedAt
	event.DeletedAt = image.DeletedAt
}

func DtoLikeImageRequestToEntityLike(ctx context.Context, req dto.LikeImageRequest, like *entity.Like) {
	userAuth := ctxuserauth.Get(ctx)
	like.UserID = userAuth.ID
	like.ImageID = req.ImageID
}

func EntityLikeToDtoImageLikedEvent(like entity.Like, event *dto.ImageLikedEvent) {
	event.ID = like.ID
	event.UserID = like.UserID
	event.ImageID = like.ImageID
	event.CreatedAt = like.CreatedAt
	event.UpdatedAt = like.UpdatedAt
	event.DeletedAt = like.DeletedAt
}

func DtoCommentImageRequestToEntityComment(ctx context.Context, req dto.CommentImageRequest, comment *entity.Comment) {
	userAuth := ctxuserauth.Get(ctx)
	comment.UserID = userAuth.ID
	comment.ImageID = req.ImageID
	comment.Comment = req.Comment
}

func EntityCommentToDtoImageCommentedEvent(comment entity.Comment, event *dto.ImageCommentedEvent) {
	event.ID = comment.ID
	event.UserID = comment.UserID
	event.ImageID = comment.ImageID
	event.Comment = comment.Comment
	event.CreatedAt = comment.CreatedAt
	event.UpdatedAt = comment.UpdatedAt
	event.DeletedAt = comment.DeletedAt
}

func EntityImageToDtoImageResponse(image entity.Image, res *dto.ImageResponse) {
	res.ID = image.ID
	res.UserID = image.UserID
	res.Caption = image.Caption
	res.URL = image.URL
	res.LikeCount = image.LikeCount
	res.CommentCount = image.CommentCount
	res.CreatedAt = image.CreatedAt
	res.UpdatedAt = image.UpdatedAt
	res.DeletedAt = image.DeletedAt
}

func EntityLikeToDtoLikeResponse(like entity.Like, res *dto.LikeResponse) {
	res.ID = like.ID
	res.UserID = like.UserID
	res.ImageID = like.ImageID
	res.CreatedAt = like.CreatedAt
	res.UpdatedAt = like.UpdatedAt
	res.DeletedAt = like.DeletedAt
}

func EntityLikeListToDtoLikeResponseList(likeList entity.LikeList, res *dto.LikeResponseList) {
	for _, like := range likeList {
		r := dto.LikeResponse{}
		EntityLikeToDtoLikeResponse(like, &r)
		*res = append(*res, r)
	}
}

func EntityCommentToDtoCommentResponse(comment entity.Comment, res *dto.CommentResponse) {
	res.ID = comment.ID
	res.UserID = comment.UserID
	res.ImageID = comment.ImageID
	res.Comment = comment.Comment
	res.CreatedAt = comment.CreatedAt
	res.UpdatedAt = comment.UpdatedAt
	res.DeletedAt = comment.DeletedAt
}

func EntityCommentListToDtoCommentResponseList(commentList entity.CommentList, res *dto.CommentResponseList) {
	for _, comment := range commentList {
		r := dto.CommentResponse{}
		EntityCommentToDtoCommentResponse(comment, &r)
		*res = append(*res, r)
	}
}

func DtoImageUploadedEventToDtoNotifyFollowerOnUploadRequest(event dto.ImageUploadedEvent, req *dto.NotifyFollowerOnUploadRequest) {
	req.UserID = event.UserID
	req.URL = event.URL
}

func DtoImageUploadedEventToDtoSyncImageToElasticsearchRequest(event dto.ImageUploadedEvent, req *dto.SyncImageToElasticsearchRequest) {
	req.ID = event.ID
	req.UserID = event.UserID
	req.Caption = event.Caption
	req.URL = event.URL
	req.LikeCount = event.LikeCount
	req.CommentCount = event.CommentCount
	req.CreatedAt = event.CreatedAt
	req.UpdatedAt = event.UpdatedAt
	req.DeletedAt = event.DeletedAt
}

func DtoSyncImageToElasticsearchRequestToDtoImageDocument(req dto.SyncImageToElasticsearchRequest, imageDocument *dto.ImageDocument) {
	imageDocument.ID = req.ID
	imageDocument.UserID = req.UserID
	imageDocument.Caption = req.Caption
	imageDocument.URL = req.URL
	imageDocument.LikeCount = req.LikeCount
	imageDocument.CommentCount = req.CommentCount
	imageDocument.CreatedAt = req.CreatedAt
	imageDocument.UpdatedAt = req.UpdatedAt
	imageDocument.DeletedAt = req.DeletedAt
}

func DtoImageCommentedEventToDtoNotifyUserImageCommentedRequest(event dto.ImageCommentedEvent, req *dto.NotifyUserImageCommentedRequest) {
	req.ImageID = event.ImageID
	req.CommenterUserID = event.UserID
}

func KGoRecordListToDtoBatchUpdateImageCommentCountRequest(ctx context.Context, records []*kgo.Record, req *dto.BatchUpdateImageCommentCountRequest) {
	mapCounter := make(map[int64]int)
	for _, record := range records {
		event := dto.ImageCommentedEvent{}
		err := json.Unmarshal(record.Value, &event)
		if err != nil {
			x.Logger.WithContext(ctx).WithError(err).Warn("Failed to unmarshal image commented event")
			continue
		}
		mapCounter[event.ImageID]++
	}

	for imageID, count := range mapCounter {
		object := dto.ImageIncreaseCommentCount{
			ImageID: imageID,
			Count:   count,
		}
		req.ImageIncreaseCommentCountList = append(req.ImageIncreaseCommentCountList, object)
	}
}

func DtoImageLikedEventToDtoNotifyUserImageLikedRequest(event dto.ImageLikedEvent, req *dto.NotifyUserImageLikedRequest) {
	req.ImageID = event.ImageID
	req.LikerUserID = event.UserID
}

func KGoRecordListToDtoBatchUpdateImageLikeCountRequest(ctx context.Context, records []*kgo.Record, req *dto.BatchUpdateImageLikeCountRequest) {
	mapCounter := make(map[int64]int)
	for _, record := range records {
		event := dto.ImageLikedEvent{}
		err := json.Unmarshal(record.Value, &event)
		if err != nil {
			x.Logger.WithContext(ctx).WithError(err).Warn("Failed to unmarshal image liked event")
			continue
		}
		mapCounter[event.ImageID]++
	}

	for imageID, count := range mapCounter {
		object := dto.ImageIncreaseLikeCount{
			ImageID: imageID,
			Count:   count,
		}
		req.ImageIncreaseLikeCountList = append(req.ImageIncreaseLikeCountList, object)
	}
}
