package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应
	config.Producer.Return.Successes = true

	// 使用给定代理地址和配置创建一个同步生产者
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}

	defer producer.Close()

	//构建发送的消息，
	msg := &sarama.ProducerMessage {
		//Topic: "test",//包含了消息的主题
		Partition: int32(10),//
		Key:        sarama.StringEncoder("key"),//
	}

	var value string
	var msgType string
	for {
		_, err := fmt.Scanf("%s", &value)
		if err != nil {
			break
		}
		fmt.Scanf("%s",&msgType)
		fmt.Println("msgType = ",msgType,",value = ",value)
		msg.Topic = msgType
		//将字符串转换为字节数组
		msg.Value = sarama.ByteEncoder(value)
		//fmt.Println(value)
		//SendMessage：该方法是生产者生产给定的消息
		//生产成功的时候返回该消息的分区和所在的偏移量
		//生产失败的时候返回error
		partition, offset, err := producer.SendMessage(msg)

		if err != nil {
			fmt.Println("Send message Fail")
		}
		fmt.Printf("Partition = %d, offset=%d\n", partition, offset)
	}
}
ctx := context.Background()
client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://127.0.0.1:9200"))
HandleError(err, "newclient")

// 用IndexExists检查索引是否存在
exists, err := client.IndexExists(indexName).Do(ctx)
HandleError(err, "indexexist")
fmt.Println("Phone No. = ")
if !exists {
	// 用CreateIndex创建索引，mapping内容用BodyString传入
	_, err := client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
	HandleError(err, "createindex")
}
fmt.Println("Phone No. =bbb ")


// 写入
docEs, err := client.Index().
	Index(indexName).
	Id(strconv.Itoa(data_res[0].ID)).
	BodyJson(data_res[0]).
	Refresh("wait_for").
	Do(ctx)

HandleError(err, "clientindex")
fmt.Printf("Indexed with id=%v, type=%s\n", docEs.Id, docEs.Type)
//读取
result, err := client.Get().
	Index(indexName).
	Id(strconv.Itoa(data_res[0].ID)).
	Do(ctx)
HandleError(err, "clientget")
if result.Found {
	fmt.Printf("Got document %v (version=%d, index=%s, type=%s)\n",
		result.Id, result.Version, result.Index, result.Type)
	err := json.Unmarshal(result.Source, &data_res_back)
	HandleError(err, "clientfound")
	fmt.Println(data_res_back.ID, data_res_back.Title, data_res_back.Source, data_res_back.Types, data_res_back.Timestamp, data_res_back.Body)
}
//大量写入
bulkRequest := client.Bulk()
for _, subject := range data_res {
	doc := elastic.NewBulkIndexRequest().Index(indexName).Id(strconv.Itoa(subject.ID)).Doc(subject)
	bulkRequest = bulkRequest.Add(doc)
}

response, err := bulkRequest.Do(ctx)
HandleError(err, "bulkrequest")
failed := response.Failed()
l := len(failed)
if l > 0 {
	fmt.Printf("Error(%d)", l, response.Errors)
}