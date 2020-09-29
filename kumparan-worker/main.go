package main

import (
	"context"
	"encoding/json"
	"fmt"
	"kumparan/kumparan-worker/db"
	"kumparan/kumparan-worker/utils"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/olivere/elastic"

	"strings"
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
	}
	`
)

const (
	defaultKafkaTopics   = "kumparan" //""pulse-telkomsel-req""
	defaultConsumerGroup = "kumparan-client"
	defaultZookeeper     = "127.0.0.1:2181"
	defaultBroker        = "127.0.0.1:9092,127.0.0.1:9093,127.0.0.1:9094"
)

var (
	consumerGroup string
	kafkaTopics   string
	zookeeper     string
	urlakuisisi   string
	brokers       string
	debughttp     bool
	elasticHost   string
	client        *elastic.Client
)

type PublishNews struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

func init() {
	consumerGroup = utils.GetEnv("KAFKA_GROUP", defaultConsumerGroup)
	brokers = utils.GetEnv("KAFKA_BROKERS", defaultBroker)
	kafkaTopics = utils.GetEnv("KAFKA_TOPICS", defaultKafkaTopics)
	zookeeper = utils.GetEnv("KAFKA_ZOOKEEPER", defaultZookeeper)
	elasticHost = utils.GetEnv("ELASTIC_HOST", "http://127.0.0.1:9200")

	var err error
	client, err = elastic.NewClient(
		elastic.SetURL(elasticHost), // http://localstack:4571
		elastic.SetSniff(false),     //Keep the list of IPs initially given, and don't care about updates in the cluster.
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "ELASTIC ", log.LstdFlags)),
		elastic.SetTraceLog(log.New(os.Stdout, "ELASTIC ", log.LstdFlags)),
	)

	if err != nil {
		fmt.Println("Found error when add client elastic: ", err)
		return
	}
}

func main() {
	fmt.Println("Worker Kafka Init Start")
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	fmt.Println("List Brokers :", brokers)

	fmt.Println("Kafka Topics :", kafkaTopics)
	listBroker := strings.Split(brokers, ",")
	topics := strings.Split(kafkaTopics, ",")
	consumer, err := cluster.NewConsumer(listBroker, consumerGroup, topics, config)

	if err != nil {
		fmt.Println("Cluster Connect problem ", err)
	}
	// consume errors
	go func() {
		for err := range consumer.Errors() {
			fmt.Println("Consumer Error: %s\n", err.Error())
			//fmt.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			fmt.Println("Rebalanced: %+v\n", ntf)
			//fmt.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				processMsg := ProcessMessage{
					Consumer: consumer,
					Message:  msg,
				}
				go processMsg.SaveNews()
				/// ///go SendingNewFds(msg.Value) ============>>>
				//go SendingFds(PackMessageReqFds(msg.Value))
			}
		}
	}
}

type ProcessMessage struct {
	Consumer *cluster.Consumer
	Message  *sarama.ConsumerMessage
}

func (process *ProcessMessage) SaveNews() {
	dataReq := PublishNews{}
	if err := json.Unmarshal(process.Message.Value, &dataReq); err != nil {
		fmt.Println("Error where unmarshal data from kafka: ", err)

		return
	}

	// save to db
	dbMsg := db.News{
		Author: dataReq.Author,
		Body:   dataReq.Body,
	}
	dbs := db.InitNewsDB()
	err := dbs.SaveNews(&dbMsg)
	if err != nil {
		fmt.Println("Failed to save news to db: ", err)

		return
	}

	// save to es
	ctx, stop := context.WithTimeout(context.Background(), time.Second)
	defer stop()

	exists, err := client.IndexExists("kumparan").Do(ctx)
	if err != nil {
		// Handle error

		fmt.Println("Error check index: ", err)
		return
	}
	if !exists {
		// Create a new index.
		_, err := client.CreateIndex("kumparan").BodyString(mappings).Do(ctx)
		if err != nil {
			// Handle error

			fmt.Println("Error create new index: ", err)
			return
		}
	}

	dataElastic := map[string]interface{}{
		"id":         dbMsg.ID,
		"created_at": time.Now().Format("2006-01-02 15:04:05"),
	}
	_, err = client.Index().Index("kumparan").Type("news").Id(strconv.Itoa(dbMsg.ID)).BodyJson(dataElastic).Do(ctx)
	if err != nil {
		fmt.Println("Failed to add data: ", err)
		return
	}

	process.Consumer.MarkOffset(process.Message, "")
	return
}
