package common
type Modify_data struct {
	ID                      int            `json:"id"`
	Timestamp               string `json:"Timestamp"`
	Source					string `json:"Source"`
	Title               string `json:"title"`
	Body        string `json:"body"`
	Types              []string `json:"Types"`

}
