package review_rating_incremented

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"
)

type ReviewRatingIncrementedEventMarshaller interface {
	Marshal(event *ReviewRatingIncrementedEventMessage) ([]byte, error)
	Unmarshal([]byte) (*ReviewRatingIncrementedEventMessage, error)
}

type ReviewRatingIncrementedEventJSONMarshaller struct {
}

func (r *ReviewRatingIncrementedEventJSONMarshaller) Marshal(event *ReviewRatingIncrementedEventMessage) ([]byte, error) {
	marshalledMessage, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return marshalledMessage, nil
}

func (r *ReviewRatingIncrementedEventJSONMarshaller) Unmarshal(bytes []byte) (*ReviewRatingIncrementedEventMessage, error) {
	reviewCreatedEvent := &ReviewRatingIncrementedEventMessage{}
	err := json.Unmarshal(bytes, reviewCreatedEvent)
	if err != nil {
		return nil, err
	}

	return reviewCreatedEvent, nil
}

type ReviewRatingIncrementedEventProtobufMarshaller struct {
}

func (r *ReviewRatingIncrementedEventProtobufMarshaller) Marshal(event *ReviewRatingIncrementedEventMessage) ([]byte, error) {
	protobufMessage := &ReviewRatingIncrementedEvent{
		Review: &ReviewRatingIncrementedEvent_Review{
			Uuid: event.Review.UUID,
		},
	}

	return proto.Marshal(protobufMessage)
}

func (r *ReviewRatingIncrementedEventProtobufMarshaller) Unmarshal(bytes []byte) (*ReviewRatingIncrementedEventMessage, error) {
	reviewCreatedEvent := &ReviewRatingIncrementedEvent{}
	err := proto.Unmarshal(bytes, reviewCreatedEvent)
	if err != nil {
		return nil, err
	}

	return &ReviewRatingIncrementedEventMessage{
		Review: ReviewMessage{
			UUID: reviewCreatedEvent.Review.GetUuid(),
		},
	}, nil
}
