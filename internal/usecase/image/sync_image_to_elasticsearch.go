package image

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) SyncImageToElasticsearch(ctx context.Context, req dto.SyncImageToElasticsearchRequest) error {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "image.(*ImageUsecaseImpl).SyncImageToElasticsearch")
	}

	imageDocument := dto.ImageDocument{}
	converter.DtoSyncImageToElasticsearchRequestToDtoImageDocument(req, &imageDocument)

	err = u.ImageSearch.IndexImage(ctx, &imageDocument)
	if err != nil {
		return errkit.AddFuncName(err, "image.(*ImageUsecaseImpl).SyncImageToElasticsearch")
	}

	return nil
}
