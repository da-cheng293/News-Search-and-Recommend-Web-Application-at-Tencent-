import (
"bytes"
"encoding/json"
"fmt"
"io/ioutil"
"net/http"
)

func main() {
fmt.Println("Starting the application...")
response, err := http.Get("https://httpbin.org/ip")
if err != nil {
fmt.Printf("The HTTP request failed with error %s\n", err)
} else {
data, _ := ioutil.ReadAll(response.Body)
fmt.Println(string(data))
}
jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
jsonValue, _ := json.Marshal(jsonData)
response, err = http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(jsonValue))
if err != nil {
fmt.Printf("The HTTP request failed with error %s\n", err)
} else {
data, _ := ioutil.ReadAll(response.Body)
fmt.Println(string(data))
}
fmt.Println("Terminating the application...")
}


package main

import (
"fmt"
"net/http"
"ideo/crawldata"
)


func myweb(w http.ResponseWriter, r *http.Request) {
fmt.Fprintf(w,"just a start")
}
func main() {

crawldata.Crawl()
http.HandleFunc("/", myweb)
fmt.Println("hello world")
err :=http.ListenAndServe(":8080",nil)
if err!=nil{
fmt.Println("wrong")
}
}

package main

import (
"bytes"
"context"
"encoding/json"
"fmt"
"github.com/PuerkitoBio/goquery"
"github.com/olivere/elastic"
"github.com/yanyiwu/gojieba"
"io/ioutil"
"log"
"net/http"
"reflect"
"regexp"
"strconv"
"strings"
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
//subject   Modify_data
indexName = "data_res"
typeName = "online"
servers   = "http://localhost:9200/"
)



var find bool   // 寻找文章标志位
//var listFind bool // 寻找文章列表标志位
var buf bytes.Buffer
//var num int = 10  //爬取章节数量
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
type Modify_data struct {
ID                      int            `json:"id"`
Timestamp               string `json:"Timestamp"`
Source					string `json:"Source"`
Title               string `json:"title"`
Body        string `json:"body"`
Types              []string `json:"Types"`


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
func Getbody(url string)  {
pageStr := GetPageStr(url)
//fmt.Println(pageStr)
re := regexp.MustCompile(`<div class="article" id="article">"([\u4e00-\u9fa5])"`)
results := re.FindAllStringSubmatch(pageStr, -1)
fmt.Printf("找到%d条结果:\n",len(results))
for _, result := range results {
//fmt.Println(result)
fmt.Println(result[1])
}
}
//func everyChapter(target string)  {
//	fmt.Println(target)
//	resp,err := http.Get(target)
//	if err!=nil{
//		fmt.Println("get err http",err)
//	}
//	defer resp.Body.Close()
//	doc,err := html.Parse(resp.Body)
//	find = false
//	parse(doc)
//解析文章
//func parse(n *html.Node)  {
//	if n.Type == html.ElementNode && n.Data == "div"{
//		for _,a := range n.Attr{
//			if a.Key == "id" && a.Val == "content" {
//				find = true
//				parseTxt(&buf,n)
//				break
//			}
//		}
//	}
//	if !find{
//		for c := n.FirstChild;c!=nil;c=c.NextSibling{
//			parse(c)
//		}
//	}
//}
//	//提取文字
//func parseTxt(buf *bytes.Buffer,n *html.Node)  {
//	for c:=n.FirstChild;c!=nil;c=c.NextSibling{
//		if c.Data != "br"{
//			buf.WriteString(c.Data+"\n")
//		}
//	}
//}
//func myweb(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w,"just a start")
//}
//func everyChapter(target string) {
//	fmt.Println(target)
//	resp, err := http.Get(target)
//	if err != nil {
//		fmt.Println("get err http", err)
//	}
//	defer resp.Body.Close()
//	doc, err := html.Parse(resp.Body)
//	find = false
//	fmt.Println("OK", doc)
//	parse(doc)
//	fmt.Println("OKfind", find)
//	text,err := os.Create("sina.txt")
//	if err!=nil{
//		fmt.Println("get create file err",err)
//	}
//	file := strings.NewReader(buf.String())
//	file.WriteTo(text)
//}
//func parse(n *html.Node) {
//	if n.Type == html.ElementNode && n.Data == "div" {
//		for _, a := range n.Attr {
//			if a.Key == "id" && a.Val == "article" {
//				find = true
//				parseTxt(&buf, n)
//				break
//			}
//		}
//	}
//	if !find{
//		for c := n.FirstChild;c!=nil;c=c.NextSibling{
//			parse(c)
//		}
//	}
//}
//func parseTxt(buf *bytes.Buffer,n *html.Node)  {
//	for c:=n.FirstChild;c!=nil;c=c.NextSibling{
//		if c.Data != "br"{
//			buf.WriteString(c.Data+"\n")
//			//fmt.Println("OK", c.Data)
//		}
//	}
//}
func main() {

var responseObject soa
//var test string
var data_res Modify_data
var data_res_back Modify_data

url := "http://api.tianapi.com/generalnews/index?key=e522570c5b2737fb6be17f0184bd87d1&page=1&&num=20"
req, _ := http.NewRequest("GET", url, nil)
res, _ := http.DefaultClient.Do(req)
defer res.Body.Close()
fmt.Println("var1 = ", reflect.TypeOf(res.Body))
body, _ := ioutil.ReadAll(res.Body)
fmt.Println("Phone No. = ", string(body))
//var responseObject soa
//var test string
json.Unmarshal(body, &responseObject)


fmt.Println(len(responseObject.Newslist))
for i := 0; i < len(responseObject.Newslist); i++ {
fmt.Println(responseObject.Newslist[i].Url)
}
doc, err := goquery.NewDocumentFromReader(strings.NewReader(GetPageStr(responseObject.Newslist[1].Url)))
fmt.Println(responseObject.Newslist[1].Url)
fmt.Println(responseObject.Newslist[1].Ctime)
if err != nil {

log.Fatal(err)

}
data_res.Title=responseObject.Newslist[1].Title
doc.Find(".article").Find("p").Each(func(i int, selection *goquery.Selection) {
//fmt.Println(selection.Text())
data_res.Body=data_res.Body+selection.Text()
})
fmt.Println(data_res.Body)
x := gojieba.NewJieba()
defer x.Free()

keywords := x.ExtractWithWeight(data_res.Body, 5)
fmt.Println("Extract:", keywords)
fmt.Println(keywords[0].Word)
for _, elem := range keywords{
data_res.Types=append(data_res.Types, elem.Word)
}
//doc.Find("#keywords").Each(func(i int, selection *goquery.Selection) {
//	fmt.Println(selection.Text())
//})
doc.Find("span.author").Each(func(i int, selection *goquery.Selection) {
fmt.Println(selection.Text())
data_res.Source=data_res.Source+selection.Text()

})
doc.Find("span.date").Each(func(i int, selection *goquery.Selection) {
fmt.Println(selection.Text())
data_res.Timestamp=data_res.Timestamp+selection.Text()
})
data_res.ID=1
ctx := context.Background()
client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://127.0.0.1:9200"))
if err != nil {
panic(err)
}

// 用IndexExists检查索引是否存在
exists, err := client.IndexExists(indexName).Do(ctx)
if err != nil {
panic(err)
}
fmt.Println("Phone No. = ")
if !exists {
// 用CreateIndex创建索引，mapping内容用BodyString传入
_, err := client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
if err != nil {
panic(err)
}
}
fmt.Println("Phone No. =bbb ")
//subject = Subject{
//	ID:     1,
//	Title:  "肖恩克的救赎",
//	Genres: []string{"犯罪", "剧情"},
//}

// 写入
docc, err := client.Index().
Index(indexName).
Id(strconv.Itoa(data_res.ID)).
BodyJson(data_res).
Refresh("wait_for").
Do(ctx)

if err != nil {
panic(err)
}
fmt.Printf("Indexed with id=%v, type=%s\n", docc.Id, docc.Type)

result, err := client.Get().
Index(indexName).
Id(strconv.Itoa(data_res.ID)).
Do(ctx)
if err != nil {
panic(err)
}
if result.Found {
fmt.Printf("Got document %v (version=%d, index=%s, type=%s)\n",
result.Id, result.Version, result.Index, result.Type)
err := json.Unmarshal(result.Source, &data_res_back)
if err != nil {
panic(err)
}
fmt.Println(data_res_back.ID, data_res_back.Title, data_res_back.Source, data_res_back.Types, data_res_back.Timestamp, data_res_back.Body)
}
//GetPageStr(responseObject.Newslist[0].Url)
//resp,err := http.Get(responseObject.Newslist[0].Url)
//if err!=nil{
//	fmt.Println("get err http",err)
//	return err
//}
//doc,err := html.Parse(resp.Body)
//if err != nil{
//	fmt.Println("html parse err",err)
//	return err
//}
//parseList(doc)
//var record []Numverify
//for body in res.Body:
//if err := json.NewDecoder(res.Body).Decode(&record); err != nil {
//log.Println(err)
//}
//json.Unmarshal([]byte(string(body)), &record)

//fmt.Println("Phone No. = ", record)

//fmt.Println(string(body))
}
//jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
//jsonValue, _ := json.Marshal(jsonData)
//response, err = http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(jsonValue))
//if err != nil {
//	fmt.Printf("The HTTP request failed with error %s\n", err)
//} else {
//	data, _ := ioutil.ReadAll(response.Body)
//	fmt.Println(string(data))
//}


