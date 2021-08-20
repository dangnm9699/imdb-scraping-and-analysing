package model

// Movie represents a movie, includes lots of information
type Movie struct {
	Url             string          `json:"url,omitempty" bson:"url"`
	Name            string          `json:"name,omitempty" bson:"name"`
	AlternateName   string          `json:"alternateName,omitempty" bson:"alternateName"`
	Image           string          `json:"image,omitempty" bson:"image"`
	ContentRating   string          `json:"contentRating,omitempty" bson:"contentRating"`
	Genre           []string        `json:"genre,omitempty" bson:"genre"`
	Actor           []Object        `json:"actor,omitempty" bson:"actor"`
	Director        []Object        `json:"director,omitempty" bson:"director"`
	Creator         []Object        `json:"creator,omitempty" bson:"creator"`
	Trailer         Trailer         `json:"trailer,omitempty" bson:"trailer"`
	DatePublished   string          `json:"datePublished,omitempty" bson:"datePublished"`
	Description     string          `json:"description,omitempty" bson:"description"`
	Keywords        string          `json:"keywords,omitempty" bson:"keywords"`
	AggregateRating AggregateRating `json:"aggregateRating,omitempty" bson:"aggregateRating"`
	Duration        string          `json:"duration,omitempty" bson:"duration"`
}

// MovieMsg represents a message to communicate in channel
type MovieMsg struct {
	Url string
	Raw string
}
