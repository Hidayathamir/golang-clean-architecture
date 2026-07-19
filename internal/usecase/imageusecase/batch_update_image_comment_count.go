package imageusecase

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

func (u *ImageUsecaseImpl) BatchUpdateImageCommentCount(ctx context.Context, req dto.BatchUpdateImageCommentCountRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).BatchUpdateImageCommentCount")
	}

	for _, v := range req.ImageIncreaseCommentCountList {
		err = u.ImageRepository.IncrementCommentCountByID(ctx, u.DB, v.ImageID, v.Count)
		if err != nil {
			logkit.Logger.WithContext(ctx).WithError(err).WithField("v", v).Warn()
		}
	}

	return nil
}
