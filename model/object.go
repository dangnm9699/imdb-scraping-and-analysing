package model

type Object struct {
	Type string `json:"@type,omitempty" bson:"type"`
	Url  string `json:"url,omitempty" bson:"url"`
	Name string `json:"name,omitempty" bson:"name"`
}
