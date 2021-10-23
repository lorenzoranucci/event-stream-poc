package kafka

type ReviewCreatedEventMessage struct {
	Review ReviewMessage `json:"review"`
}

type ReviewMessage struct {
	UUID    string `json:"uuid"`
	Comment string `json:"comment"`
	Rating  int32  `json:"rating"`
}
