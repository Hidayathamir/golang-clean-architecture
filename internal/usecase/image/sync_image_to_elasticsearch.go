package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) SyncImageToElasticsearch(ctx context.Context, req dto.SyncImageToElasticsearchRequest) error {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	imageDocument := dto.ImageDocument{}
	converter.DtoSyncImageToElasticsearchRequestToDtoImageDocument(req, &imageDocument)

	err = u.ImageSearch.IndexImage(ctx, &imageDocument)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
