package model

type Trailer struct {
	Type         string `json:"@type,omitempty"`
	Name         string `json:"name,omitempty"`
	EmbedUrl     string `json:"embedUrl,omitempty"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
	Description  string `json:"description,omitempty"`
}
