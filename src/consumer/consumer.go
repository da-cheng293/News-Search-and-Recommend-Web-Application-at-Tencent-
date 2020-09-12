package main

import (

	"common"

)

//type common.Modify_data struct {
//	ID                      int            `json:"id"`
//	Timestamp               string `json:"Timestamp"`
//	Source					string `json:"Source"`
//	Title               string `json:"title"`
//	Body        string `json:"body"`
//	Types              []string `json:"Types"`
//
//}

type Consumer struct {
	ready chan bool
}

//var data_res []common.Modify_data
//
//func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
//	close(consumer.ready)
//	return nil
//}
//
//func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
//	return nil
//}
//
//func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim ) error {
//	for message := range claim.Messages() {
//		err := msgpack.Unmarshal(message.Value, &data_res)
//		if err != nil {
//			panic(err)
//		}
//
//
//		log.Printf("Message claimed: key = %s, Subject(id=%d, Timestamp,=%s, title=%s, body=%s, type=%v, source=%s )", string(message.Key),
//			data_res[0].ID, data_res[0].Timestamp, data_res[0].Title, data_res[0].Body,data_res[0].Types,data_res[0].Source)
//		session.MarkMessage(message, "")
//	}
//
//	return nil
//}

func main() {
//kafaka 消费数据
	broker := "localhost:9092"
	group := "1"
	var topics []string
	topics=append(topics, "message_pack")
	common.Consumer_pro(broker, group, topics)
	//version, err := sarama.ParseKafkaVersion("2.3.0")
	//if err != nil {
	//	log.Panicf("Error parsing Kafka version: %v", err)
	//}
	//
	//config := sarama.NewConfig()
	//config.Version = version
	//consumer := Consumer{
	//	ready: make(chan bool, 0),
	//}
	//
	//ctx, cancel := context.WithCancel(context.Background())
	//client, err := sarama.NewConsumerGroup([]string{broker}, group, config)
	//if err != nil {
	//	log.Panicf("Error creating consumer group client: %v", err)
	//}
	//
	//wg := &sync.WaitGroup{}
	//go func() {
	//	wg.Add(1)
	//	defer wg.Done()
	//	for {
	//		if err := client.Consume(ctx, topics, &consumer); err != nil {
	//			log.Panicf("Error from consumer: %v", err)
	//		}
	//		if ctx.Err() != nil {
	//			return
	//		}
	//		fmt.Println("okkk")
	//		consumer.ready = make(chan bool, 0)
	//	}
	//}()
	//
	//<-consumer.ready
	//
	//sigterm := make(chan os.Signal, 1)
	//signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	//select {
	//case <-sigterm:
	//	log.Println("terminating: via signal")
	//}
	//cancel()
	//wg.Wait()
	//if err = client.Close(); err != nil {
	//	log.Panicf("Error closing client: %v", err)
	//}
}