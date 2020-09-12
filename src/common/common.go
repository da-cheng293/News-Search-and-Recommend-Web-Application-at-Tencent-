package common

import "fmt"

type Modify_data struct {
	ID                      int            `json:"id"`
	Timestamp               string `json:"Timestamp"`
	Source					string `json:"Source"`
	Title               string `json:"title"`
	Body        string `json:"body"`
	Url_News         string  `json:"url"`
	Types              []string `json:"Types"`

}

type Consumer struct {
	ready chan bool
}


func HandleError(err error, why string) {
	if err != nil {
		fmt.Println(why, err)
	}
}