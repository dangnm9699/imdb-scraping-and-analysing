package model

type AggregateRating struct {
	Type        string  `json:"@type,omitempty" bson:"type"`
	RatingCount int     `json:"ratingCount,omitempty" bson:"ratingCount"`
	BestRating  int     `json:"bestRating,omitempty" bson:"bestRating"`
	WorstRating int     `json:"worstRating,omitempty" bson:"worstRating"`
	RatingValue float64 `json:"ratingValue,omitempty" bson:"ratingValue"`
}
