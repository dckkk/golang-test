package services

import (
	"errors"
	"fmt"

	"kumparan/kumparan-api/kafka"

	"github.com/Shopify/sarama"
)

type ServicePush struct {
	Topic       string
	Brokerlist  string
	ConfigKafka *sarama.Config
}

// Push is function for publish to kafka
func (s *ServicePush) Push(bdata []byte) error {
	//fmt.Sprintf("Data : %s",string(bdata))
	kafkaProducer := kafka.GetConectionKafka(s.Brokerlist, s.ConfigKafka)
	if kafkaProducer.ErrRes != nil {
		fmt.Println("Err Client REQ :", kafkaProducer.ErrRes)
		return kafkaProducer.ErrRes
	}
	//msqreq := model.GenerateMessageQuee(req)
	kafkaProducer.Topic = s.Topic //fmt.Sprintf("%s-req", s.Topic)
	kafkaProducer.Value = sarama.ByteEncoder(bdata)
	kafkaProducer.Process()
	if kafkaProducer.Done == false {
		return errors.New("Sending to Kafka Failed")
	}
	//go s.Db.SavetransactionLog(msqreq)
	return nil
}

// CheckConnectionKafka
func (s *ServicePush) CheckConnectionKafka() (bool, error) {
	kafkaProducer := kafka.GetConectionKafka(s.Brokerlist, s.ConfigKafka)
	if kafkaProducer.ErrRes != nil {
		return false, kafkaProducer.ErrRes
	}
	return true, kafkaProducer.ErrRes
}
