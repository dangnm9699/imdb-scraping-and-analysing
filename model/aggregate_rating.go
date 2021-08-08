package model

type AggregateRating struct {
	Type        string  `json:"@type,omitempty"`
	RatingCount int     `json:"ratingCount,omitempty"`
	BestRating  int     `json:"bestRating,omitempty"`
	WorstRating int     `json:"worstRating,omitempty"`
	RatingValue float64 `json:"ratingValue,omitempty"`
}
