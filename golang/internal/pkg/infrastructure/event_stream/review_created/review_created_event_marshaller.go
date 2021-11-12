package review_created

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"
)

type ReviewCreatedEventMarshaller interface {
	Marshal(event *ReviewCreatedEventMessage) ([]byte, error)
	Unmarshal([]byte) (*ReviewCreatedEventMessage, error)
}

type ReviewCreatedEventJSONMarshaller struct {
}

func (r *ReviewCreatedEventJSONMarshaller) Marshal(event *ReviewCreatedEventMessage) ([]byte, error) {
	marshalledMessage, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return marshalledMessage, nil
}

func (r *ReviewCreatedEventJSONMarshaller) Unmarshal(bytes []byte) (*ReviewCreatedEventMessage, error) {
	reviewCreatedEvent := &ReviewCreatedEventMessage{}
	err := json.Unmarshal(bytes, reviewCreatedEvent)
	if err != nil {
		return nil, err
	}

	return reviewCreatedEvent, nil
}

type ReviewCreatedEventProtobufMarshaller struct {
}

func (r *ReviewCreatedEventProtobufMarshaller) Marshal(event *ReviewCreatedEventMessage) ([]byte, error) {
	protobufMessage := &ReviewCreatedEvent{
		Review: &ReviewCreatedEvent_Review{
			Uuid:    event.Review.UUID,
			Comment: event.Review.Comment,
			Rating:  event.Review.Rating,
		},
	}

	return proto.Marshal(protobufMessage)
}

func (r *ReviewCreatedEventProtobufMarshaller) Unmarshal(bytes []byte) (*ReviewCreatedEventMessage, error) {
	reviewCreatedEvent := &ReviewCreatedEvent{}
	err := proto.Unmarshal(bytes, reviewCreatedEvent)
	if err != nil {
		return nil, err
	}

	return &ReviewCreatedEventMessage{
		Review: ReviewMessage{
			UUID:    reviewCreatedEvent.Review.GetUuid(),
			Comment: reviewCreatedEvent.Review.GetComment(),
			Rating:  reviewCreatedEvent.Review.GetRating(),
		},
	}, nil
}
