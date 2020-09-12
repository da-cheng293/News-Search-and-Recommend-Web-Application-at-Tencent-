package main

import (
	"common"
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
			"Types": {
				"type": "text"
			}
		}
	}
}`

var (

	indexName = "data_res"
	typeName = "online"
	servers   = "http://localhost:9200/"
)


type soa struct {
	Newslist    []Newslist `json:"Newslist"`
}


type Newslist struct {

	Ctime               string `json:"ctime"`
	Title               string `json:"title"`
	Description         string `json:"description"`
	PicUrl              string `json:"picUrl"`
	Url                 string `json:"url"`

}

//type common.Modify_data struct {
//	ID                      int            `json:"id"`
//	Timestamp               string `json:"Timestamp"`
//	Source					string `json:"Source"`
//	Title               string `json:"title"`
//	Body        string `json:"body"`
//	Types              []string `json:"Types"`
//
//}
type fn func (common.Modify_data, string) common.Modify_data
func newschina_deal(single_data_res common.Modify_data, news_china_url string) common.Modify_data{
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(GetPageStr(news_china_url)))

	HandleError(err, "goquery")


	doc.Find("#chan_newsDetail").Find("p").Each(func(i int, selection *goquery.Selection) {

		single_data_res.Body=single_data_res.Body+selection.Text()
	})
	fmt.Println(single_data_res.Body)
	x := gojieba.NewJieba()
	defer x.Free()

	keywords := x.ExtractWithWeight(single_data_res.Body, 5)
	fmt.Println("Extract:", keywords)

	for _, elem := range keywords{
		single_data_res.Types=append(single_data_res.Types, elem.Word)
	}

	doc.Find("span.source").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
		single_data_res.Source=single_data_res.Source+selection.Text()

	})


	//single_data_res.ID=sub_id
	return single_data_res
}
func sina_deal(single_data_res common.Modify_data, sina_url string) common.Modify_data{
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(GetPageStr(sina_url)))

		HandleError(err, "goquery")


		doc.Find(".article").Find("p").Each(func(i int, selection *goquery.Selection) {

			single_data_res.Body=single_data_res.Body+selection.Text()
		})
		fmt.Println(single_data_res.Body)
		x := gojieba.NewJieba()
		defer x.Free()

		keywords := x.ExtractWithWeight(single_data_res.Body, 5)
		fmt.Println("Extract:", keywords)

		for _, elem := range keywords{
			single_data_res.Types=append(single_data_res.Types, elem.Word)
		}

		doc.Find("span.author").Each(func(i int, selection *goquery.Selection) {
			fmt.Println(selection.Text())
			single_data_res.Source=single_data_res.Source+selection.Text()

		})


		//single_data_res.ID=sub_id
		return single_data_res
}
func zhibo_deal(single_data_res common.Modify_data, zhibo_url string) common.Modify_data{
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(GetPageStr(zhibo_url)))

	HandleError(err, "goquery")


	doc.Find(".article-content").Find("p").Each(func(i int, selection *goquery.Selection) {

		single_data_res.Body=single_data_res.Body+selection.Text()
	})
	fmt.Println(single_data_res.Body)
	x := gojieba.NewJieba()
	defer x.Free()

	keywords := x.ExtractWithWeight(single_data_res.Body, 5)
	fmt.Println("Extract:", keywords)

	for _, elem := range keywords{
		single_data_res.Types=append(single_data_res.Types, elem.Word)
	}

	doc.Find("span.anchor").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
		single_data_res.Source=single_data_res.Source+selection.Text()

	})


	//single_data_res.ID=sub_id
	return single_data_res
}
func bar(msg string) {
	log.Printf("bar! Message is %s", msg)
}

func HandleError(err error, why string) {
	if err != nil {
		fmt.Println(why, err)
	}
}

func GetPageStr(url string) (pageStr string) {
	//1.发送http请求，获取页面内容
	resp, err := http.Get(url)
	//处理异常
	HandleError(err, "http.Get url")
	//关闭资源
	defer resp.Body.Close()
	//接收页面
	pageBytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "ioutil.ReadAll")
	//打印页面内容
	pageStr = string(pageBytes)
	//fmt.Println("OK", pageStr)
	return pageStr
}


func Get_datares(data *[]common.Modify_data, sub_title string, sub_time string, sub_url string, sub_id int)  {

	var single_data_res common.Modify_data


	parseUrl, err := url.Parse(sub_url)
	HandleError(err, "parseUrl")

	fmt.Println(parseUrl.Host)
	HandleError(err, "parseUrl.Host")

	single_data_res.ID=sub_id
	single_data_res.Title=sub_title
	single_data_res.Timestamp=sub_time
	//定制爬取程序
	url_map := map[string] fn {
		"k.sina.com.cn": sina_deal,
		"news.sina.com.cn": sina_deal,
		"news.china.com": newschina_deal,
		"v.zhibo.tv": zhibo_deal,
	}

	value, ok := url_map[parseUrl.Host]
	if ok {
		*data = append(*data, value(single_data_res, sub_url))
	} else {
		fmt.Println("key not found")
	}

}


func main() {

	var responseObject soa
	var data_res []common.Modify_data   //十个新闻数据 的arrary
	//var data_res_back common.Modify_data  //接受 从es 读取的数据
	var id_sub int



	urlApi := "http://api.tianapi.com/generalnews/index?key=e522570c5b2737fb6be17f0184bd87d1&page=1&&num=10"
	req, _ := http.NewRequest("GET", urlApi, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	fmt.Println("var1 = ", reflect.TypeOf(res.Body))
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("Phone No. = ", string(body))

	json.Unmarshal(body, &responseObject)


	fmt.Println(len(responseObject.Newslist))

	for i := 0; i < len(responseObject.Newslist); i++ {
		fmt.Println(responseObject.Newslist[i].Url)
		title_sub := responseObject.Newslist[i].Title
		Ctime_sub := responseObject.Newslist[i].Ctime
		url_sub := responseObject.Newslist[i].Url
		id_sub=id_sub+1
		Get_datares(&data_res, title_sub, Ctime_sub, url_sub, id_sub)
	}
	for j :=0; j<len(data_res);j++{
		fmt.Println(data_res[j])
	}

	//kafka 生产数据
	broker := "localhost:9092"
	topic := "message_pack"

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



