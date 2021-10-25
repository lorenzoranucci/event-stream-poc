package cmd

import (
	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/event_stream"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/create_review"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/in_memory"
	"github.com/ProntoPro/event-stream-golang/pkg/kafka"
	"github.com/ProntoPro/event-stream-golang/pkg/pulsar"
)

const (
	KafkaEventStream  = "kafka"
	PulsarEventStream = "pulsar"
	JSONFormat        = "json"
	ProtobufFormat    = "protobuf"
)

func newServiceLocator(kafkaURL string, pulsarURL string, eventStream string, format string) *serviceLocator {
	return &serviceLocator{
		kafkaURL:        kafkaURL,
		pulsarURL:       pulsarURL,
		eventStreamKind: eventStream,
		format:          format,
	}
}

type serviceLocator struct {
	createReviewCommandHandler *application.CreateReviewCommandHandler

	createReviewHandler  *create_review.CreateReviewHandler
	createReviewConsumer event_stream.Consumer
	createReviewProducer event_stream.Producer

	reviewCreatedMarshaller    event_stream.ReviewCreatedEventMarshaller
	reviewCreatedEventProducer *event_stream.ReviewCreatedEventProducer
	reviewCreatedEventConsumer *event_stream.ReviewCreatedEventConsumer

	jsonMarshaller     *event_stream.ReviewCreatedEventJSONMarshaller
	protobufMarshaller *event_stream.ReviewCreatedEventProtobufMarshaller

	inMemoryReviewRepository *in_memory.ReviewRepository

	kafkaConsumer *kafka.Consumer
	kafkaProducer *kafka.Producer

	pulsarConsumer *pulsar.Consumer
	pulsarProducer *pulsar.Producer

	kafkaURL        string
	pulsarURL       string
	eventStreamKind string
	format          string
}

func (s *serviceLocator) CreateReviewCommandHandler() *application.CreateReviewCommandHandler {
	if s.createReviewCommandHandler == nil {
		s.createReviewCommandHandler = application.NewCreateReviewCommandHandler(s.ReviewRepository(), s.EventBus())
	}

	return s.createReviewCommandHandler
}

func (s *serviceLocator) ReviewRepository() application.ReviewRepository {
	return s.InMemoryReviewRepository()
}

func (s *serviceLocator) EventBus() application.EventBus {
	return s.ReviewCreatedEventProducer()
}

func (s *serviceLocator) CreateReviewHandler() *create_review.CreateReviewHandler {
	if s.createReviewHandler == nil {
		s.createReviewHandler = create_review.NewCreateReviewHandler(s.CreateReviewCommandHandler())
	}

	return s.createReviewHandler
}

func (s *serviceLocator) CreateReviewConsumer() event_stream.Consumer {
	if s.createReviewConsumer == nil {
		switch s.eventStreamKind {
		case KafkaEventStream:
			s.createReviewConsumer = s.KafkaConsumer()
		case PulsarEventStream:
			s.createReviewConsumer = s.PulsarConsumer()
		}
	}

	return s.createReviewConsumer
}

func (s *serviceLocator) CreateReviewProducer() event_stream.Producer {
	if s.createReviewProducer == nil {
		switch s.eventStreamKind {
		case KafkaEventStream:
			s.createReviewProducer = s.KafkaProducer()
		case PulsarEventStream:
			s.createReviewProducer = s.PulsarProducer()
		}
	}

	return s.createReviewProducer
}

func (s *serviceLocator) ReviewCreatedMarshaller() event_stream.ReviewCreatedEventMarshaller {
	if s.reviewCreatedMarshaller == nil {
		switch s.format {
		case JSONFormat:
			s.reviewCreatedMarshaller = s.JsonMarshaller()
		case ProtobufFormat:
			s.reviewCreatedMarshaller = s.ProtobufMarshaller()
		}
	}

	return s.reviewCreatedMarshaller
}

func (s *serviceLocator) InMemoryReviewRepository() *in_memory.ReviewRepository {
	if s.inMemoryReviewRepository == nil {
		s.inMemoryReviewRepository = &in_memory.ReviewRepository{}
	}

	return s.inMemoryReviewRepository
}

func (s *serviceLocator) ReviewCreatedEventProducer() *event_stream.ReviewCreatedEventProducer {
	if s.reviewCreatedEventProducer == nil {
		s.reviewCreatedEventProducer = event_stream.NewReviewCreatedEventProducer(
			s.CreateReviewProducer(),
			s.ReviewCreatedMarshaller(),
		)
	}

	return s.reviewCreatedEventProducer
}

func (s *serviceLocator) ReviewCreatedEventConsumer() *event_stream.ReviewCreatedEventConsumer {
	if s.reviewCreatedEventConsumer == nil {
		s.reviewCreatedEventConsumer = event_stream.NewReviewCreatedEventConsumer(
			s.CreateReviewConsumer(),
			s.ReviewCreatedMarshaller(),
		)
	}

	return s.reviewCreatedEventConsumer
}

func (s *serviceLocator) KafkaConsumer() *kafka.Consumer {
	if s.kafkaConsumer == nil {
		var err error
		s.kafkaConsumer, err = kafka.NewConsumer(s.KafkaURL(), 0, "review_created")
		if err != nil {
			panic(err)
		}
	}
	return s.kafkaConsumer
}

func (s *serviceLocator) KafkaProducer() *kafka.Producer {
	if s.kafkaProducer == nil {
		var err error
		s.kafkaProducer, err = kafka.NewProducer(s.KafkaURL(), "review_created")
		if err != nil {
			panic(err)
		}
	}
	return s.kafkaProducer
}

func (s *serviceLocator) PulsarConsumer() *pulsar.Consumer {
	if s.pulsarConsumer == nil {
		var err error
		s.pulsarConsumer, err = pulsar.NewConsumer(s.PulsarURL(), "review_created")
		if err != nil {
			panic(err)
		}
	}
	return s.pulsarConsumer
}

func (s *serviceLocator) PulsarProducer() *pulsar.Producer {
	if s.pulsarProducer == nil {
		var err error
		s.pulsarProducer, err = pulsar.NewProducer(s.PulsarURL(), "review_created")
		if err != nil {
			panic(err)
		}
	}
	return s.pulsarProducer
}

func (s *serviceLocator) JsonMarshaller() *event_stream.ReviewCreatedEventJSONMarshaller {
	if s.jsonMarshaller == nil {
		s.jsonMarshaller = &event_stream.ReviewCreatedEventJSONMarshaller{}
	}

	return s.jsonMarshaller
}

func (s *serviceLocator) ProtobufMarshaller() *event_stream.ReviewCreatedEventProtobufMarshaller {
	if s.protobufMarshaller == nil {
		s.protobufMarshaller = &event_stream.ReviewCreatedEventProtobufMarshaller{}
	}

	return s.protobufMarshaller
}

func (s *serviceLocator) KafkaURL() string {
	return s.kafkaURL
}

func (s *serviceLocator) PulsarURL() string {
	return s.pulsarURL
}

func (s *serviceLocator) EventStreamKind() string {
	return s.eventStreamKind
}

func (s *serviceLocator) Format() string {
	return s.format
}
