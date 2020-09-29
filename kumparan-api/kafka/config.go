package kafka

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

//Get Connection Kafka
func GetConectionKafka(brokerList string, config *sarama.Config) KafkaProducer {
	producer, err := sarama.NewSyncProducer(strings.Split(brokerList, ","), config)
	return KafkaProducer{
		Connection: producer,
		ErrRes:     err,
	}
}

func GetConfigKafka(clientID string, maxmsgbyte int) *sarama.Config {

	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	configKafka := sarama.NewConfig()
	configKafka.Producer.Return.Successes = true
	configKafka.Producer.Partitioner = sarama.NewRandomPartitioner
	configKafka.ClientID = clientID
	configKafka.Producer.MaxMessageBytes = maxmsgbyte
	configKafka.Producer.Timeout = 5 * time.Second
	configKafka.Net.DialTimeout = 2 * time.Second
	configKafka.Net.ReadTimeout = 5 * time.Second
	configKafka.Net.WriteTimeout = 5 * time.Second

	return configKafka
}
