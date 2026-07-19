package provider

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticsearchClient(cfg *config.Config) *elasticsearch.Client {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:  []string{cfg.GetElasticsearchAddress()},
		MaxRetries: 3,
	})
	errkit.PanicIfErr(err)

	return client
}
