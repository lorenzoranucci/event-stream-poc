package cmd

import (
	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/create_review"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/in_memory"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/kafka"
	kafka2 "github.com/ProntoPro/event-stream-golang/pkg/kafka"
)

func NewServiceLocator(kafkaURL string) (*ServiceLocator, error) {
	kafkaProducer, err := kafka2.NewProducer(kafkaURL, "review_created_event")
	if err != nil {
		return nil, err
	}

	repository := &in_memory.ReviewRepository{}

	jsonMarshaller := &kafka.ReviewCreatedEventJSONMarshaller{}
	protobufMarshaller := &kafka.ReviewCreatedEventProtobufMarshaller{}

	createReviewHandlerWithKafkaAndJSON := create_review.NewCreateReviewHandler(
		application.NewCreateReviewCommandHandler(
			repository,
			kafka.NewReviewCreatedEventProducer(
				kafkaProducer,
				jsonMarshaller,
			),
		),
	)

	createReviewHandlerWithKafkaAndProtobuf := create_review.NewCreateReviewHandler(
		application.NewCreateReviewCommandHandler(
			repository,
			kafka.NewReviewCreatedEventProducer(
				kafkaProducer,
				protobufMarshaller,
			),
		),
	)

	kafkaConsumer, err := kafka2.NewConsumer(kafkaURL, 0, "review_created_event")
	if err != nil {
		return nil, err
	}

	reviewCreatedConsumerWithKafkaAndJSON := kafka.NewReviewCreatedEventConsumer(kafkaConsumer, jsonMarshaller)

	reviewCreatedConsumerWithKafkaAndProtobuf := kafka.NewReviewCreatedEventConsumer(kafkaConsumer, protobufMarshaller)

	return &ServiceLocator{
		createReviewHandlerWithKafkaAndJSON:       createReviewHandlerWithKafkaAndJSON,
		createReviewHandlerWithKafkaAndProtobuf:   createReviewHandlerWithKafkaAndProtobuf,
		reviewCreatedConsumerWithKafkaAndJSON:     reviewCreatedConsumerWithKafkaAndJSON,
		reviewCreatedConsumerWithKafkaAndProtobuf: reviewCreatedConsumerWithKafkaAndProtobuf,
	}, nil
}

type ServiceLocator struct {
	createReviewHandlerWithKafkaAndJSON     *create_review.CreateReviewHandler
	createReviewHandlerWithKafkaAndProtobuf *create_review.CreateReviewHandler

	reviewCreatedConsumerWithKafkaAndJSON     *kafka.ReviewCreatedEventConsumer
	reviewCreatedConsumerWithKafkaAndProtobuf *kafka.ReviewCreatedEventConsumer
}

func (s *ServiceLocator) CreateReviewHandlerWithKafkaAndJSON() *create_review.CreateReviewHandler {
	return s.createReviewHandlerWithKafkaAndJSON
}

func (s *ServiceLocator) CreateReviewHandlerWithKafkaAndProtobuf() *create_review.CreateReviewHandler {
	return s.createReviewHandlerWithKafkaAndProtobuf
}

func (s *ServiceLocator) ReviewCreatedConsumerWithKafkaAndJSON() *kafka.ReviewCreatedEventConsumer {
	return s.reviewCreatedConsumerWithKafkaAndJSON
}

func (s *ServiceLocator) ReviewCreatedConsumerWithKafkaAndProtobuf() *kafka.ReviewCreatedEventConsumer {
	return s.reviewCreatedConsumerWithKafkaAndProtobuf
}
