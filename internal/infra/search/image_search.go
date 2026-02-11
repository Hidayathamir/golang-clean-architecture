package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/indexname"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/elastic/go-elasticsearch/v8"
)

//go:generate moq -out=../../mock/MockSearchImage.go -pkg=mock . ImageSearch

type ImageSearch interface {
	IndexImage(ctx context.Context, document *dto.ImageDocument) error
}

type ImageSearchImpl struct {
	client *elasticsearch.Client
}

var _ ImageSearch = &ImageSearchImpl{}

func NewImageSearch(client *elasticsearch.Client) ImageSearch {
	return &ImageSearchImpl{
		client: client,
	}
}

func (i *ImageSearchImpl) IndexImage(ctx context.Context, document *dto.ImageDocument) error {
	jsonByte, err := json.Marshal(document)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	res, err := i.client.Index(indexname.Images, bytes.NewReader(jsonByte))
	if err != nil {
		return errkit.AddFuncName(err)
	}
	defer x.LogIfErrForDeferContext(ctx, res.Body.Close)

	if res.IsError() {
		err := errors.New(res.String())
		err = errkit.Wrap(err, "indexing error")
		return errkit.AddFuncName(err)
	}

	return nil
}
