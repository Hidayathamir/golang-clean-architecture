package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/indexname"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/elastic/go-elasticsearch/v8"
)

//go:generate moq -out=../../mock/MockSearchImage.go -pkg=mock . ImageSearch

type ImageSearch interface {
	IndexImage(ctx context.Context, document *dto.ImageDocument) error
	SearchImages(ctx context.Context, imageList *entity.ImageList, cursorCreatedAt time.Time, cursorID int64, limit int) error
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
		return errkit.AddFuncName(err, "search.(*ImageSearchImpl).IndexImage")
	}

	res, err := i.client.Index(indexname.Images, bytes.NewReader(jsonByte))
	if err != nil {
		return errkit.AddFuncName(err, "search.(*ImageSearchImpl).IndexImage")
	}
	defer logkit.LogIfErrForDeferContext(ctx, res.Body.Close)

	if res.IsError() {
		err := errors.New(res.String())
		err = errkit.Wrap(err, "indexing error")
		return errkit.AddFuncName(err, "search.(*ImageSearchImpl).IndexImage")
	}

	return nil
}

// searchHit is one Elasticsearch hit wrapper; only _source is needed.
type searchHit struct {
	Source dto.ImageDocument `json:"_source"`
}

type searchResponse struct {
	Hits struct {
		Hits []searchHit `json:"hits"`
	} `json:"hits"`
}

func (i *ImageSearchImpl) SearchImages(ctx context.Context, imageList *entity.ImageList, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	// Discovery order mirrors the list endpoints: newest first, tie-break by id.
	// The keyset predicate keeps pages stable across concurrent inserts.
	sortClause := []any{
		map[string]any{"created_at": map[string]any{"order": "desc"}},
		map[string]any{"id": map[string]any{"order": "desc"}},
	}
	body := map[string]any{
		"size": limit,
		"sort": sortClause,
	}

	if cursorID > 0 {
		body["query"] = map[string]any{
			"bool": map[string]any{
				"should": []any{
					map[string]any{"range": map[string]any{"created_at": map[string]any{"lt": cursorCreatedAt}}},
					map[string]any{
						"bool": map[string]any{
							"must": []any{
								map[string]any{"term": map[string]any{"created_at": cursorCreatedAt}},
								map[string]any{"range": map[string]any{"id": map[string]any{"lt": cursorID}}},
							},
						},
					},
				},
				"minimum_should_match": 1,
			},
		}
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		return errkit.AddFuncName(err, "search.(*ImageSearchImpl).SearchImages")
	}

	res, err := i.client.Search(
		i.client.Search.WithContext(ctx),
		i.client.Search.WithIndex(indexname.Images),
		i.client.Search.WithBody(bytes.NewReader(bodyByte)),
	)
	if err != nil {
		return errkit.AddFuncName(err, "search.(*ImageSearchImpl).SearchImages")
	}
	defer logkit.LogIfErrForDeferContext(ctx, res.Body.Close)

	if res.IsError() {
		err := errors.New(res.String())
		err = errkit.Wrap(err, "search error")
		return errkit.AddFuncName(err, "search.(*ImageSearchImpl).SearchImages")
	}

	var parsed searchResponse
	err = json.NewDecoder(res.Body).Decode(&parsed)
	if err != nil {
		return errkit.AddFuncName(err, "search.(*ImageSearchImpl).SearchImages")
	}

	for _, hit := range parsed.Hits.Hits {
		image := entity.Image{
			ID:           hit.Source.ID,
			UserID:       hit.Source.UserID,
			Caption:      hit.Source.Caption,
			URL:          hit.Source.URL,
			LikeCount:    hit.Source.LikeCount,
			CommentCount: hit.Source.CommentCount,
			CreatedAt:    hit.Source.CreatedAt,
			UpdatedAt:    hit.Source.UpdatedAt,
			DeletedAt:    hit.Source.DeletedAt,
		}
		*imageList = append(*imageList, image)
	}

	return nil
}
