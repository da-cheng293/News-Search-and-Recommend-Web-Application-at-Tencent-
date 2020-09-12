package common

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/vmihailenco/msgpack"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"sync"
)

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


		log.Printf("Message claimed: key = %s, Subject(id=%d, Timestamp,=%s, title=%s, body=%s, type=%v, source=%s )", string(message.Key),
			data_res[0].ID, data_res[0].Timestamp, data_res[0].Title, data_res[0].Body,data_res[0].Types,data_res[0].Source)
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