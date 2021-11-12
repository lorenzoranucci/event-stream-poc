package review_rating_incremented

type ReviewRatingIncrementedEventMessage struct {
	Review ReviewMessage `json:"review"`
}

type ReviewMessage struct {
	UUID string `json:"uuid"`
}
