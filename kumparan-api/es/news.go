package es

import (
	"context"

	"github.com/olivere/elastic"
)

type NewsElastic interface {
	GetLatestNews() (*elastic.SearchResult, error)
}

func InitNewsElastic() NewsElastic {
	return &NewsElasticStruct{ES: client}
}

type NewsElasticStruct struct {
	ES *elastic.Client
}

func (s *NewsElasticStruct) GetLatestNews() (*elastic.SearchResult, error) {
	termQuery := elastic.NewTermQuery("id", "created_at")
	searchResult, err := s.ES.Search().
		Index("kumparan"). // search in index "twitter"
		Query(termQuery).  // specify the query
		Sort("id", false). // sort by "user" field, ascending
		From(0).Size(10).  // take documents 0-9
		Pretty(true).      // pretty print request and response JSON
		Do(context.Background())
	return searchResult, err
}
