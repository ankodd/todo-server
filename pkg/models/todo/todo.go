package todo

type Todo struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
	Done bool   `json:"done,omitempty"`
}
