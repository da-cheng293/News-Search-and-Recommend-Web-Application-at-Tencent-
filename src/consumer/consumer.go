package main

import (
	"common"
)



func main() {
//kafaka 消费数据
	broker := "localhost:9092"
	group := "1"
	var topics []string
	topics=append(topics, "message_pack")

	common.Consumer_pro(broker, group, topics)

}