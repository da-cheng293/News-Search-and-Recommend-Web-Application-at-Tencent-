package common

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/olivere/elastic"
	"github.com/vmihailenco/msgpack"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"github.com/gomodule/redigo/redis"
	"sync"
	"crypto/md5"
	"encoding/json"
	"github.com/kai1987/go-text-censor"
)
var (
	//subject   Modify_data
	IndexName = "data_res"
	typeName  = "online"
	servers   = "http://localhost:9200/"
)
const mapping = `
{
	"mappings": {

		"properties": {
			"Timestamp": {
				"type": "text"
			},
            "Source": {
				"type": "text"
			},
			"Title": {
				"type": "text"
			},
            "Body": {
				"type": "text"
			},
			"Url_News": {
				"type": "text"
			},
			"Types": {
				"type": "text"
			}
		}
	}
}`
var data_res []Modify_data
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim ) error {
	for message := range claim.Messages() {
		err := msgpack.Unmarshal(message.Value, &data_res)
		if err != nil {
			panic(err)
		}

		redisclient, err := redis.Dial("tcp", "localhost:6379")
		if err != nil {
			panic(err)
		}
		for j :=0; j<len(data_res);j++{
			fmt.Println(data_res[j])
			temp_value, _ := json.Marshal(data_res[j])
			temp_md5 := md5.Sum(temp_value)
			key := fmt.Sprintf("%x", temp_md5)
			if _, err := redisclient.Do("AUTH", "g3838438"); err != nil {
				panic(err)
			}
			n,err := redisclient.Do("SETNX", key, "")
			if err != nil {
				panic(err)
			}
			if n == int64(0){
				data_res = append(data_res[:j], data_res[j+1:]...)
			}
		}
		//文字审核
		if len(data_res)>0{
			for i :=0; i<len(data_res);i++{
				isPass := textcensor.IsPass(data_res[i].Body,true)
				if isPass == false{
					data_res = append(data_res[:i], data_res[i+1:]...)
				}

			}
		}
		//es 写入
		ctx := context.Background()
		client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://127.0.0.1:9200"))
		HandleError(err, "newclient")

		// 用IndexExists检查索引是否存在
		exists, err := client.IndexExists(IndexName).Do(ctx)
		HandleError(err, "indexexist")
		fmt.Println("Phone No. = ")
		if !exists {
			// 用CreateIndex创建索引，mapping内容用BodyString传入
			_, err := client.CreateIndex(IndexName).BodyString(mapping).Do(ctx)
			HandleError(err, "createindex")
		}
		fmt.Println("Phone No. =bbb ")

		bulkRequest := client.Bulk()
		for _, subject := range data_res {
			doc := elastic.NewBulkIndexRequest().Index(IndexName).Id(strconv.Itoa(subject.ID)).Doc(subject)
			bulkRequest = bulkRequest.Add(doc)
		}

		response, err := bulkRequest.Do(ctx)
		HandleError(err, "bulkrequest")
		failed := response.Failed()
		l := len(failed)
		if l > 0 {
			fmt.Printf("Error(%d)", l, response.Errors)
		}


		log.Printf("Message claimed: key = %s, Subject(id=%d, Timestamp,=%s, title=%s, body=%s, type=%v, source=%s, url=%s )", string(message.Key),
			data_res[0].ID, data_res[0].Timestamp, data_res[0].Title, data_res[0].Body,data_res[0].Types,data_res[0].Source, data_res[0].Url_News)
		session.MarkMessage(message, "")
	}

	return nil
}

func Produce(topic string, broker string, data_res *[]Modify_data)  {
	//kafka 生产数据


	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	//sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	producer, err := sarama.NewAsyncProducer([]string{broker}, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var enqueued, errors int
	doneCh := make(chan struct{})
	var data_id int
	data_id=1
	go func() {
		for {

			time.Sleep(1 * time.Second)

			b, err := msgpack.Marshal(&data_res)
			if err != nil {
				panic(err)
			}

			msg := &sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.StringEncoder(strconv.Itoa(data_id)),
				Value: sarama.StringEncoder(b),
			}
			select {
			case producer.Input() <- msg:
				enqueued++
				fmt.Printf("Produce message: %v\n", data_res)
			case err := <-producer.Errors():
				errors++
				fmt.Println("Failed to produce message:", err)
			case <-signals:
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)

}


func Consumer_pro(broker string, group string, topics []string){

	version, err := sarama.ParseKafkaVersion("2.3.0")
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	config := sarama.NewConfig()
	config.Version = version
	consumer := Consumer{
		ready: make(chan bool, 0),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup([]string{broker}, group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			if err := client.Consume(ctx, topics, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			fmt.Println("okkk")
			consumer.ready = make(chan bool, 0)
		}
	}()

	<-consumer.ready

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}