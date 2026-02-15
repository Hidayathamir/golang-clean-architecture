package provider

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticsearchClient(cfg *config.Config) *elasticsearch.Client {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:  []string{cfg.GetElasticsearchAddress()},
		MaxRetries: 3,
	})
	x.PanicIfErr(err)

	return client
}
