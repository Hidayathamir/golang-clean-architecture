package imageusecase

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) BatchUpdateImageLikeCount(ctx context.Context, req dto.BatchUpdateImageLikeCountRequest) error {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).BatchUpdateImageLikeCount")
	}

	for _, v := range req.ImageIncreaseLikeCountList {
		err = u.ImageRepository.IncrementLikeCountByID(ctx, u.DB, v.ImageID, v.Count)
		if err != nil {
			x.Logger.WithContext(ctx).WithError(err).WithField("v", v).Warn()
		}
	}

	return nil
}
