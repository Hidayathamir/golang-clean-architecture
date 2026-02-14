package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ ImageUsecase = &ImageUsecaseMwLogger{}

type ImageUsecaseMwLogger struct {
	Next ImageUsecase
}

func NewImageUsecaseMwLogger(next ImageUsecase) *ImageUsecaseMwLogger {
	return &ImageUsecaseMwLogger{
		Next: next,
	}
}

func (u *ImageUsecaseMwLogger) Comment(ctx context.Context, req dto.CommentImageRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.Comment(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (u *ImageUsecaseMwLogger) Like(ctx context.Context, req dto.LikeImageRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.Like(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (u *ImageUsecaseMwLogger) Upload(ctx context.Context, req dto.UploadImageRequest) (dto.ImageResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Upload(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}

func (u *ImageUsecaseMwLogger) GetImage(ctx context.Context, req dto.GetImageRequest) (dto.ImageResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.GetImage(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}

func (u *ImageUsecaseMwLogger) GetComment(ctx context.Context, req dto.GetCommentRequest) (dto.CommentResponseList, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.GetComment(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}

func (u *ImageUsecaseMwLogger) GetLike(ctx context.Context, req dto.GetLikeRequest) (dto.LikeResponseList, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.GetLike(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}

func (u *ImageUsecaseMwLogger) NotifyFollowerOnUpload(ctx context.Context, req dto.NotifyFollowerOnUploadRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.NotifyFollowerOnUpload(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (u *ImageUsecaseMwLogger) SyncImageToElasticsearch(ctx context.Context, req dto.SyncImageToElasticsearchRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.SyncImageToElasticsearch(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (u *ImageUsecaseMwLogger) NotifyUserImageCommented(ctx context.Context, req dto.NotifyUserImageCommentedRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.NotifyUserImageCommented(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (u *ImageUsecaseMwLogger) BatchUpdateImageCommentCount(ctx context.Context, req dto.BatchUpdateImageCommentCountRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.BatchUpdateImageCommentCount(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (u *ImageUsecaseMwLogger) NotifyUserImageLiked(ctx context.Context, req dto.NotifyUserImageLikedRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.NotifyUserImageLiked(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (u *ImageUsecaseMwLogger) BatchUpdateImageLikeCount(ctx context.Context, req dto.BatchUpdateImageLikeCountRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.BatchUpdateImageLikeCount(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err)

	return err
}
