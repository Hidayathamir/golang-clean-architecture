package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) BatchUpdateImageCommentCount(ctx context.Context, req *dto.BatchUpdateImageCommentCountRequest) error {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	for _, v := range req.ImageIncreaseCommentCountList {
		if err := u.ImageRepository.IncrementCommentCountByID(ctx, u.DB, v.ImageID, v.Count); err != nil {
			x.Logger.WithContext(ctx).WithError(err).WithField("v", v).Warn()
		}
	}

	return nil
}
