package model

type Trailer struct {
	Type         string `json:"@type,omitempty" bson:"type"`
	Name         string `json:"name,omitempty" bson:"name"`
	EmbedUrl     string `json:"embedUrl,omitempty" bson:"embedUrl"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty" bson:"thumbnailUrl"`
	Description  string `json:"description,omitempty" bson:"description"`
}
