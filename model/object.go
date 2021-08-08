package model

type Object struct {
	Type string `json:"@type,omitempty"`
	Url  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}
