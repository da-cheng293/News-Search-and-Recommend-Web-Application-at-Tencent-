package main
import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Shopify/sarama"
	"github.com/vmihailenco/msgpack"
	"github.com/yanyiwu/gojieba"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 去重
redisclient, err := redis.Dial("tcp", "localhost:6379")
if err != nil {
panic(err)
}
for j :=0; j<len(data_res);j++{
fmt.Println(data_res[j])
temp_value, _ := json.Marshal(data_res[j])
temp_md5 := md5.Sum(temp_value)
key := fmt.Sprintf("%x", temp_md5)
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
exists, err := client.IndexExists(indexName).Do(ctx)
HandleError(err, "indexexist")
fmt.Println("Phone No. = ")
if !exists {
// 用CreateIndex创建索引，mapping内容用BodyString传入
_, err := client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
HandleError(err, "createindex")
}
fmt.Println("Phone No. =bbb ")

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