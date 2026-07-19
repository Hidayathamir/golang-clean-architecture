package imageusecase

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
	"gorm.io/gorm"
)

func (u *ImageUsecaseImpl) Comment(ctx context.Context, req dto.CommentImageRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Comment")
	}

	comment := entity.Comment{}
	converter.DtoCommentImageRequestToEntityComment(ctx, req, &comment)

	event := dto.ImageCommentedEvent{}
	converter.EntityCommentToDtoImageCommentedEvent(comment, &event)

	err = u.DB.Transaction(func(tx *gorm.DB) error {
		err := u.CommentRepository.Create(ctx, tx, &comment)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Comment")
		}

		err = u.ImageProducer.SendImageCommented(ctx, tx, &event)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Comment")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
