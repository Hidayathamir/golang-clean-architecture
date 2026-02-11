package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) SyncImageToElasticsearch(ctx context.Context, req *dto.SyncImageToElasticsearchRequest) error {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	imageDocument := new(dto.ImageDocument)
	converter.DtoSyncImageToElasticsearchRequestToDtoImageDocument(ctx, req, imageDocument)

	if err := u.ImageSearch.IndexImage(ctx, imageDocument); err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
