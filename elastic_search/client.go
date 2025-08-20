package elasticsearch

import (
	"log"

	"github.com/Matltin/Bilitioo-Backend/util"
	"github.com/elastic/go-elasticsearch"
)

func NewElasticsearchClient(config util.Config) *elasticsearch.Client {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.ElasticsearchAddress},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return es
}