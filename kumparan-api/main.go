package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"kumparan/kumparan-api/db"
	"kumparan/kumparan-api/es"
	"kumparan/kumparan-api/kafka"
	kafkaservices "kumparan/kumparan-api/kafka/services"
	"kumparan/kumparan-api/models"
	"kumparan/kumparan-api/utils"
)

func init() {
	kafkaBroker := utils.GetEnv("KAFKA_BROKERS", "127.0.0.1:9092,127.0.0.1:9093,127.0.0.1:9094")
	topicNews := utils.GetEnv("KUMPARAN_TOPIC", "kumparan-news")
	kafkaConfig := kafka.GetConfigKafka("kumparan", 50000000)
	kafkaServices := kafkaservices.ServicePush{
		Topic:       topicNews,
		Brokerlist:  kafkaBroker,
		ConfigKafka: kafkaConfig,
	}
	kafkaServices.CheckConnectionKafka()
}

func main() {
	listenAddress := utils.GetEnv("APP_ADDRESS", ":7878")
	http.HandleFunc("/news", News)

	fmt.Println("starting web server at: ", listenAddress)
	http.ListenAndServe(listenAddress, nil)
}

// News is func to handle route with endpoint /news
func News(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		PublishNews(w, r)
	case "GET":
		GetNews(w, r)
	default:
	}
}

// GetNews is func to get 10 latest news
func GetNews(w http.ResponseWriter, r *http.Request) {
	res := models.ResponseNews{
		Code:    "01",
		Message: "Failed",
	}

	elastic := es.InitNewsElastic()
	searchResult, err := elastic.GetLatestNews()
	if err != nil {
		fmt.Println("Failed to search data to elastic: ", err)

		result, _ := json.Marshal(res)
		w.Write(result)
		return
	}

	resData := []db.News{}
	// Start goroutine
	if searchResult.TotalHits() > 0 {
		wg := sync.WaitGroup{}
		wg.Add(int(searchResult.TotalHits()))
		esResultChan := make(chan es.News)
		for _, hit := range searchResult.Hits.Hits {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				var esResult es.News
				err := json.Unmarshal(*hit.Source, &esResult)
				if err != nil {
					// Deserialization failed
					fmt.Println("Failed to unmarshal data from es: ", err)
					return
				}
				esResultChan <- esResult
			}(&wg)
		}

		go func() {
			wg.Wait()
			close(esResultChan)
		}()

		for item := range esResultChan {
			dbs := db.InitNewsDB()
			dbRes, err := dbs.GetDetailNews(item.ID)
			if err != nil {
				fmt.Println("Failed to get detail news from db: ", err)
				continue
			}
			resData = append(resData, dbRes)
		}
		//  == End Goroutine
	}

	res.Code = "00"
	res.Message = "Succesful"
	res.Data = resData
	result, _ := json.Marshal(res)
	w.Write(result)
	return
}

// PublishNews is func to handle publish news api
func PublishNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req := models.RequestPublishNews{}
	res := models.ResponseNews{
		Code:    "01",
		Message: "Failed",
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		fmt.Println("Bad request")

		result, _ := json.Marshal(res)
		w.Write(result)
		return
	}

	kafkaBroker := utils.GetEnv("KAFKA_BROKERS", "127.0.0.1:9092,127.0.0.1:9093,127.0.0.1:9094")
	topicNews := utils.GetEnv("KUMPARAN_TOPIC", "kumparan-news")
	kafkaConfig := kafka.GetConfigKafka("kumparan", 50000000)
	kafkaServices := kafkaservices.ServicePush{
		Topic:       topicNews,
		Brokerlist:  kafkaBroker,
		ConfigKafka: kafkaConfig,
	}
	// kafkaServices.CheckConnectionKafka()
	jsonReq, _ := json.Marshal(req)
	err := kafkaServices.Push(jsonReq)
	if err != nil {
		fmt.Println("Failed send data to kafka")

		result, _ := json.Marshal(res)
		w.Write(result)
		return
	}

	res.Code = "00"
	res.Message = "Succesful"
	result, _ := json.Marshal(res)
	w.Write(result)
	return
}
