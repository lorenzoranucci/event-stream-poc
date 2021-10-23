package kafka

type ReviewCreatedEventMessage struct {
	Review ReviewMessage `json:"review"`
}

type ReviewMessage struct {
	UUID    string `json:"uuid"`
	Comment string `json:"comment"`
	Rating  int    `json:"rating"`
}
