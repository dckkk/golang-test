package es

import (
	"fmt"
	"log"
	"os"

	"kumparan/kumparan-api/utils"

	"github.com/olivere/elastic"
)

var (
	elasticHost string
	client      *elastic.Client
)

const (
	mappings = `
	{
		"settings":{
			"number_of_shards":2,
			"number_of_replicas":1
		},
		"mappings":{
			"properties":{
				"field str":{
					"type":"text"
				},
				"field int":{
					"type":"integer"
				},
				"field bool":{
					"type":"boolean"
				}
			}
		}
	}`
)

func init() {
	elasticHost = utils.GetEnv("ELASTIC_HOST", "http://127.0.0.1:9200")

	var err error
	client, err = elastic.NewClient(
		elastic.SetURL(elasticHost),
		elastic.SetSniff(false),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "ELASTIC ", log.LstdFlags)),
		elastic.SetTraceLog(log.New(os.Stdout, "ELASTIC ", log.LstdFlags)),
	)

	if err != nil {
		fmt.Println("Found error when add client elastic: ", err)
		return
	}
}
