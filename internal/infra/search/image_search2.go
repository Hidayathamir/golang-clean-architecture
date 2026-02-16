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
	"github.com/opensearch-project/opensearch-go/v2"
)

//go:generate moq -out=../../mock/MockSearchImage2.go -pkg=mock . ImageSearch2

type ImageSearch2 interface {
	IndexImage(ctx context.Context, document *dto.ImageDocument) error
}

type ImageSearch2Impl struct {
	client *opensearch.Client
}

var _ ImageSearch2 = &ImageSearch2Impl{}

func NewImageSearch2(client *opensearch.Client) ImageSearch2 {
	return &ImageSearch2Impl{
		client: client,
	}
}

func (i *ImageSearch2Impl) IndexImage(ctx context.Context, document *dto.ImageDocument) error {
	jsonByte, err := json.Marshal(document)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	res, err := i.client.Index(
		indexname.Images,
		bytes.NewReader(jsonByte),
		i.client.Index.WithContext(ctx),
	)
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
