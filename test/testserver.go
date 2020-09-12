
package main

import (
	"fmt"
	"log"
	"net/http"
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
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}


	fmt.Fprintf(w, "Hello!")
}
func main() {
	//kafaka 消费数据
	go func() {
		http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintf(w, `405`)
				return
			}
			fmt.Fprintf(w, `pong`)
		})
		http.ListenAndServe("localhost:8082", nil)
	}()

	http.HandleFunc("/statinfo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `405`)
			return
		}
		//laklahalfhajfkalgdhf
		fmt.Fprintf(w, "op")
	})
	http.ListenAndServe("localhost:8081", nil)


	//common.Consumer_pro(broker, group, topics)
	//http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	//	if r.Method != "GET" {
	//		w.WriteHeader(http.StatusMethodNotAllowed)
	//		fmt.Fprintf(w, `405`)
	//		return
	//	}
	//	fmt.Fprintf(w, `pong`)
	//})
	//http.ListenAndServe(":80", nil)

	//func helloHandler(w http.ResponseWriter, r *http.Request) {
	//	if r.URL.Path != "/hello" {
	//		http.Error(w, "404 not found.", http.StatusNotFound)
	//		return
	//	}
	//
	//	if r.Method != "GET" {
	//		http.Error(w, "Method is not supported.", http.StatusNotFound)
	//		return
	//	}
	//
	//
	//	fmt.Fprintf(w, "Hello!")
	//}


	//func main() {
	http.HandleFunc("/hello", helloHandler) // Update this line of code


	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("ok")
	//}
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
