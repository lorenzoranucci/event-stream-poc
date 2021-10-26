package cmd

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/event_stream"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/create_review"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/get_reviews"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/in_memory"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/mysql"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/projectors"
	"github.com/ProntoPro/event-stream-golang/pkg/kafka"
	"github.com/ProntoPro/event-stream-golang/pkg/pulsar"
)

const (
	KafkaEventStream  = "kafka"
	PulsarEventStream = "pulsar"
	JSONFormat        = "json"
	ProtobufFormat    = "protobuf"
)

func newServiceLocator(
	kafkaURL string,
	pulsarURL string,
	mysqlURL string,
	eventStream string,
	format string,
) *serviceLocator {
	return &serviceLocator{
		kafkaURL:        kafkaURL,
		pulsarURL:       pulsarURL,
		mysqlURL:        mysqlURL,
		eventStreamKind: eventStream,
		format:          format,
	}
}

type serviceLocator struct {
	createReviewCommandHandler *application.CreateReviewCommandHandler
	getReviewsQueryHandler     *application.GetReviewsQueryHandler

	createReviewHandler  *create_review.CreateReviewHandler
	getReviewsHandler    *get_reviews.GetReviewsHandler
	createReviewConsumer event_stream.Consumer
	createReviewProducer event_stream.Producer

	reviewCreatedMarshaller    event_stream.ReviewCreatedEventMarshaller
	reviewCreatedEventProducer *event_stream.ReviewCreatedEventProducer
	reviewCreatedEventConsumer *event_stream.ReviewCreatedEventConsumer

	jsonMarshaller     *event_stream.ReviewCreatedEventJSONMarshaller
	protobufMarshaller *event_stream.ReviewCreatedEventProtobufMarshaller

	inMemoryCreateReviewRepository *in_memory.CreateReviewRepository
	inMemoryGetReviewsRepository   *in_memory.GetReviewsRepository

	mysqlCreateReviewRepository *mysql.CreateReviewRepository
	mysqlGetReviewsRepository   *mysql.GetReviewsRepository

	getReviewsProjector *projectors.GetReviewsProjector

	kafkaConsumer *kafka.Consumer
	kafkaProducer *kafka.Producer

	pulsarConsumer *pulsar.Consumer
	pulsarProducer *pulsar.Producer

	mysqlDB *sqlx.DB

	kafkaURL        string
	pulsarURL       string
	mysqlURL        string
	eventStreamKind string
	format          string
}

func (s *serviceLocator) MysqlCreateReviewRepository() *mysql.CreateReviewRepository {
	if s.mysqlCreateReviewRepository == nil {
		s.mysqlCreateReviewRepository = mysql.NewCreateReviewRepository(s.MysqlDB())
	}
	return s.mysqlCreateReviewRepository
}

func (s *serviceLocator) MysqlGetReviewsRepository() *mysql.GetReviewsRepository {
	if s.mysqlGetReviewsRepository == nil {
		s.mysqlGetReviewsRepository = mysql.NewGetReviewsRepository(s.MysqlDB())
	}
	return s.mysqlGetReviewsRepository
}

func (s *serviceLocator) GetReviewsProjector() *projectors.GetReviewsProjector {
	if s.getReviewsProjector == nil {
		s.getReviewsProjector = projectors.NewGetReviewsProjector(s.MysqlGetReviewsRepository())
	}
	return s.getReviewsProjector
}

func (s *serviceLocator) CreateReviewCommandHandler() *application.CreateReviewCommandHandler {
	if s.createReviewCommandHandler == nil {
		s.createReviewCommandHandler = application.NewCreateReviewCommandHandler(s.CreateReviewRepository(), s.EventBus())
	}

	return s.createReviewCommandHandler
}

func (s *serviceLocator) GetReviewsQueryHandler() *application.GetReviewsQueryHandler {
	if s.getReviewsQueryHandler == nil {
		s.getReviewsQueryHandler = application.NewGetReviewsQueryHandler(s.GetReviewsRepository())
	}

	return s.getReviewsQueryHandler
}

func (s *serviceLocator) CreateReviewRepository() application.CreateReviewRepository {
	return s.MysqlCreateReviewRepository()
}

func (s *serviceLocator) GetReviewsRepository() application.GetReviewsRepository {
	return s.MysqlGetReviewsRepository()
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

func (s *serviceLocator) GetReviewsHandler() *get_reviews.GetReviewsHandler {
	if s.getReviewsHandler == nil {
		s.getReviewsHandler = get_reviews.NewGetReviewsHandler(s.GetReviewsQueryHandler())
	}

	return s.getReviewsHandler
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

func (s *serviceLocator) InMemoryCreateReviewRepository() *in_memory.CreateReviewRepository {
	if s.inMemoryCreateReviewRepository == nil {
		s.inMemoryCreateReviewRepository = &in_memory.CreateReviewRepository{}
	}

	return s.inMemoryCreateReviewRepository
}

func (s *serviceLocator) InMemoryGetReviewsRepository() *in_memory.GetReviewsRepository {
	if s.inMemoryGetReviewsRepository == nil {
		s.inMemoryGetReviewsRepository = &in_memory.GetReviewsRepository{}
	}

	return s.inMemoryGetReviewsRepository
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
			s.GetReviewsProjector(),
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

func (s *serviceLocator) MysqlDB() *sqlx.DB {
	if s.mysqlDB == nil {
		var err error
		s.mysqlDB, err = sqlx.Connect("mysql", s.MysqlURL())
		if err != nil {
			panic(err)
		}
	}
	return s.mysqlDB
}

func (s *serviceLocator) MysqlURL() string {
	return s.mysqlURL
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
