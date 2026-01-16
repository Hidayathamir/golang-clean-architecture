package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
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

func (u *ImageUsecaseMwLogger) Comment(ctx context.Context, req *model.CommentImageRequest) error {
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

func (u *ImageUsecaseMwLogger) Like(ctx context.Context, req *model.LikeImageRequest) error {
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

func (u *ImageUsecaseMwLogger) Upload(ctx context.Context, req *model.UploadImageRequest) (*model.ImageResponse, error) {
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
