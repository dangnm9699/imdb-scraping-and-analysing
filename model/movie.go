package model

type Movie struct {
	Tconst         string `json:"tconst" bson:"tconst"`
	Name           string `json:"name" bson:"name"`
	ReleasedYear   string `json:"released_year" bson:"released_year"`
	Rating         string `json:"rating" bson:"rating"`
	RatingCount    string `json:"rating_count" bson:"rating_count"`
	Runtime        string `json:"runtime" bson:"runtime"`
	Genres         string `json:"genres" bson:"genres"`
	Budget         string `json:"budget" bson:"budget"`
	GrossWorldwide string `json:"gross_worldwide" bson:"gross_worldwide"`
	Director       string `json:"director" bson:"director"`
	Stars          string `json:"stars" bson:"stars"`
	Country        string `json:"country" bson:"country"`
	StoryLine      string `json:"story_line" bson:"story_line"`
}
