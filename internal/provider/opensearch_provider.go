package provider

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/opensearch-project/opensearch-go/v2"
)

func NewOpensearchClient(cfg *config.Config) *opensearch.Client {
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses:  []string{cfg.GetOpensearchAddress()},
		MaxRetries: 3,
	})
	x.PanicIfErr(err)

	return client
}
