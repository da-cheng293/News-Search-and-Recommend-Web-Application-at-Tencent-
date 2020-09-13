package common

// Elastic search client
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"log"
	"reflect"
	"strconv"
	"time"

)



var client *elastic.Client

// NewElasticSearchClient returns an elastic seach client
func NewElasticSearchClient() *elastic.Client {
	var err error
	connected := false
	retries := 0

	// Custom retry strategy for docker-compose initialization
	for connected == false {
		// Create a new elastic client
		client, err = elastic.NewClient(
			elastic.SetURL("http://127.0.0.1:9200"), elastic.SetSniff(false))
		if err != nil {
			// log.Fatal(err)
			if retries == 5 {
				log.Fatal(err)
			}
			fmt.Println("Elasticsearch isn't ready for connection", 5-retries, "less")
			retries++
			time.Sleep(3 * time.Second)
		} else {
			connected = true
		}
	}

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	return client
}

// ExistsIndex checks if the given index exists or not
func ExistsIndex(i string) bool {
	// Check if index exists
	exists, err := client.IndexExists(i).Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return exists
}

// CreateIndex creates a new index
func CreateIndex(i string) {
	createIndex, err := client.CreateIndex(IndexName).
		Body(mapping).
		Do(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}
	if !createIndex.Acknowledged {
		log.Println("CreateIndex was not acknowledged. Check that timeout value is correct.")
	}
}

// SearchContent returns the results for a given query
func SearchContent(input string) []Modify_data {
	pages := []Modify_data{}

	ctx := context.Background()
	// Search for a page in the database using multi match query
	q := elastic.NewMultiMatchQuery(input, "Source", "Title", "Body", "Types").
		Type("most_fields").
		Fuzziness("2")
	result, err := client.Search().
		Index(IndexName).
		Query(q).
		From(0).Size(20).
		Sort("_score", false).
		Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var ttyp Modify_data
	for _, page := range result.Each(reflect.TypeOf(ttyp)) {
		p := page.(Modify_data)
		pages = append(pages, p)
	}

	return pages
}

func ReadEs()Modify_data{
	var subject Modify_data
	ctx := context.Background()
	result, err := client.Get().
		Index(IndexName).
		Id(strconv.Itoa(1)).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	if result.Found {

		err := json.Unmarshal(result.Source, &subject)
		if err != nil {
			panic(err)
		}
		

	}
	return subject

}