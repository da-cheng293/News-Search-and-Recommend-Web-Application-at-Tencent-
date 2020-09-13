package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"common"
	"github.com/gorilla/mux"
)



// SearchResult struct to handle search queries
type SearchResult struct {
	Pages []common.Modify_data `json:"pages"`
	Input string `json:"input"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/views/home.html")
	if err != nil {
		log.Print("Template parsing error: ", err)
	}
	if r.Method == "POST" {
		ajax_post_data := r.FormValue("ajax_post_data")
		fmt.Println("Receive ajax post data string ", ajax_post_data)
		w.Write([]byte("<h2>after<h2>"))
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Print("Template executing error: ", err)
	}
	//result_sin :=common.ReadEs()
	//
	//jsonData, err := json.Marshal(result_sin)
	//if err != nil {
	//	log.
	//		Print("JSON executing error: ", err)
	//	return
	//}
	//
	//w.Header().Set("Content-Type", "application/json")
	//w.Write(jsonData)
	////if r.Method == "POST" {
	////	ajax_post_data := r.FormValue("ajax_post_data")
	////	fmt.Println("Receive ajax post data string ", ajax_post_data)
	////	w.Write([]byte("<h2>after<h2>"))
	////}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	searchInput := r.Form.Get("input")

	log.Print("Querying database for: ", searchInput)

	pages := common.SearchContent(searchInput)

	searchResult := SearchResult{
		Input: searchInput,
		Pages: pages,
	}

	jsonData, err := json.Marshal(searchResult)
	if err != nil {
		log.
			Print("JSON executing error: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	common.NewElasticSearchClient()
	exists := common.ExistsIndex(common.IndexName)

	if !exists {
		common.CreateIndex(common.IndexName)
	}

	mux := mux.NewRouter()

	mux.HandleFunc("/", homeHandler).Methods("GET")
	mux.HandleFunc("/search", searchHandler).Methods("GET")
	http.ListenAndServe(":8083", mux)
}