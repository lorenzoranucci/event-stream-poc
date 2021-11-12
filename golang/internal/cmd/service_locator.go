package cmd

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/event_stream/review_rating_incremented"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/patch_review"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/commands"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/queries"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/event_stream/review_created"
	get_reviews2 "github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/projectors/get_reviews"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/event_stream"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/create_review"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/get_reviews"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/mysql"
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
	createReviewCommandHandler          *commands.CreateReviewCommandHandler
	getReviewsQueryHandler              *queries.GetReviewsQueryHandler
	incrementReviewRatingCommandHandler *commands.IncrementReviewRatingCommandHandler

	createReviewHandler           *create_review.CreateReviewHandler
	getReviewsHandler             *get_reviews.GetReviewsHandler
	patchReviewHandler            *patch_review.PatchReviewHandler
	createReviewConsumer          event_stream.Consumer
	createReviewProducer          event_stream.Producer
	incrementReviewRatingConsumer event_stream.Consumer
	incrementReviewRatingProducer event_stream.Producer

	reviewCreatedMarshaller              review_created.ReviewCreatedEventMarshaller
	reviewCreatedEventProducer           *review_created.ReviewCreatedEventProducer
	reviewCreatedEventConsumer           *review_created.ReviewCreatedEventConsumer
	reviewCreatedEventJSONMarshaller     *review_created.ReviewCreatedEventJSONMarshaller
	reviewCreatedEventProtobufMarshaller *review_created.ReviewCreatedEventProtobufMarshaller

	reviewRatingIncrementedMarshaller              review_rating_incremented.ReviewRatingIncrementedEventMarshaller
	reviewRatingIncrementedEventProducer           *review_rating_incremented.ReviewRatingIncrementedEventProducer
	reviewRatingIncrementedEventConsumer           *review_rating_incremented.ReviewRatingIncrementedEventConsumer
	reviewRatingIncrementedEventJSONMarshaller     *review_rating_incremented.ReviewRatingIncrementedEventJSONMarshaller
	reviewRatingIncrementedEventProtobufMarshaller *review_rating_incremented.ReviewRatingIncrementedEventProtobufMarshaller

	mysqlCommandsReviewRepository         *mysql.CommandsReviewRepository
	mysqlQueriesReviewRepository          *mysql.QueriesReviewsRepository
	mysqlIntegrationEventOutboxRepository *mysql.IntegrationReviewEventsOutboxRepository
	mysqlTransactionManager               *mysql.TransactionManager

	reviewCreatedProjector           *get_reviews2.ReviewCreatedProjector
	reviewRatingIncrementedProjector *get_reviews2.ReviewRatingIncrementedProjector

	reviewCreatedKafkaConsumer *kafka.Consumer
	reviewCreatedKafkaProducer *kafka.Producer

	reviewCreatedPulsarConsumer *pulsar.Consumer
	reviewCreatedPulsarProducer *pulsar.Producer

	reviewRatingIncrementedKafkaConsumer *kafka.Consumer
	reviewRatingIncrementedKafkaProducer *kafka.Producer

	reviewRatingIncrementedPulsarConsumer *pulsar.Consumer
	reviewRatingIncrementedPulsarProducer *pulsar.Producer

	mysqlDB *sqlx.DB

	kafkaURL        string
	pulsarURL       string
	mysqlURL        string
	eventStreamKind string
	format          string
}

func (s *serviceLocator) GetReviewsProjector() *get_reviews2.ReviewCreatedProjector {
	if s.reviewCreatedProjector == nil {
		s.reviewCreatedProjector = get_reviews2.NewReviewCreatedProjector(s.MysqlQueryReviewsRepository())
	}
	return s.reviewCreatedProjector
}

func (s *serviceLocator) ReviewRatingIncrementedProjector() *get_reviews2.ReviewRatingIncrementedProjector {
	if s.reviewRatingIncrementedProjector == nil {
		s.reviewRatingIncrementedProjector = get_reviews2.NewReviewRatingIncrementedProjector(
			s.MysqlQueryReviewsRepository(),
		)
	}
	return s.reviewRatingIncrementedProjector
}

func (s *serviceLocator) CreateReviewCommandHandler() *commands.CreateReviewCommandHandler {
	if s.createReviewCommandHandler == nil {
		s.createReviewCommandHandler = commands.NewCreateReviewCommandHandler(
			s.CommandReviewRepository(),
			s.MysqlTransactionManager(),
			s.MysqlIntegrationEventOutboxRepository(),
			s.ReviewCreatedEventBus(),
		)
	}

	return s.createReviewCommandHandler
}

func (s *serviceLocator) GetReviewsQueryHandler() *queries.GetReviewsQueryHandler {
	if s.getReviewsQueryHandler == nil {
		s.getReviewsQueryHandler = queries.NewGetReviewsQueryHandler(s.QueryReviewRepository())
	}

	return s.getReviewsQueryHandler
}

func (s *serviceLocator) IncrementReviewRatingCommandHandler() *commands.IncrementReviewRatingCommandHandler {
	if s.incrementReviewRatingCommandHandler == nil {
		s.incrementReviewRatingCommandHandler = commands.NewIncrementReviewRatingCommandHandler(
			s.CommandReviewRepository(),
			s.MysqlTransactionManager(),
			s.MysqlIntegrationEventOutboxRepository(),
			s.ReviewCreatedEventBus(),
		)
	}

	return s.incrementReviewRatingCommandHandler
}

func (s *serviceLocator) CommandReviewRepository() commands.ReviewRepository {
	return s.MysqlCommandReviewRepository()
}

func (s *serviceLocator) QueryReviewRepository() queries.ReviewRepository {
	return s.MysqlQueryReviewsRepository()
}

func (s *serviceLocator) ReviewCreatedEventBus() commands.IntegrationEventBus {
	return s.ReviewCreatedEventProducer()
}

func (s *serviceLocator) ReviewRatingIncrementedEventBus() commands.IntegrationEventBus {
	return s.ReviewRatingIncrementedEventProducer()
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

func (s *serviceLocator) PatchReviewHandler() *patch_review.PatchReviewHandler {
	if s.patchReviewHandler == nil {
		s.patchReviewHandler = patch_review.NewPatchReviewHandler(s.IncrementReviewRatingCommandHandler())
	}

	return s.patchReviewHandler
}

func (s *serviceLocator) CreateReviewConsumer() event_stream.Consumer {
	if s.createReviewConsumer == nil {
		switch s.eventStreamKind {
		case KafkaEventStream:
			s.createReviewConsumer = s.ReviewCreatedKafkaConsumer()
		case PulsarEventStream:
			s.createReviewConsumer = s.ReviewCreatedPulsarConsumer()
		}
	}

	return s.createReviewConsumer
}

func (s *serviceLocator) CreateReviewProducer() event_stream.Producer {
	if s.createReviewProducer == nil {
		switch s.eventStreamKind {
		case KafkaEventStream:
			s.createReviewProducer = s.ReviewCreatedKafkaProducer()
		case PulsarEventStream:
			s.createReviewProducer = s.ReviewCreatedPulsarProducer()
		}
	}

	return s.createReviewProducer
}

func (s *serviceLocator) IncrementReviewRatingConsumer() event_stream.Consumer {
	if s.incrementReviewRatingConsumer == nil {
		switch s.eventStreamKind {
		case KafkaEventStream:
			s.incrementReviewRatingConsumer = s.ReviewRatingIncrementedKafkaConsumer()
		case PulsarEventStream:
			s.incrementReviewRatingConsumer = s.ReviewRatingIncrementedPulsarConsumer()
		}
	}

	return s.incrementReviewRatingConsumer
}

func (s *serviceLocator) IncrementReviewRatingProducer() event_stream.Producer {
	if s.incrementReviewRatingProducer == nil {
		switch s.eventStreamKind {
		case KafkaEventStream:
			s.incrementReviewRatingProducer = s.ReviewCreatedKafkaProducer()
		case PulsarEventStream:
			s.incrementReviewRatingProducer = s.ReviewCreatedPulsarProducer()
		}
	}

	return s.incrementReviewRatingProducer
}

func (s *serviceLocator) ReviewCreatedMarshaller() review_created.ReviewCreatedEventMarshaller {
	if s.reviewCreatedMarshaller == nil {
		switch s.format {
		case JSONFormat:
			s.reviewCreatedMarshaller = s.ReviewCreatedEventJSONMarshaller()
		case ProtobufFormat:
			s.reviewCreatedMarshaller = s.ReviewCreatedEventProtobufMarshaller()
		}
	}

	return s.reviewCreatedMarshaller
}

func (s *serviceLocator) ReviewRatingIncrementedMarshaller() review_rating_incremented.ReviewRatingIncrementedEventMarshaller {
	if s.reviewRatingIncrementedMarshaller == nil {
		switch s.format {
		case JSONFormat:
			s.reviewRatingIncrementedMarshaller = s.ReviewRatingIncrementedEventJSONMarshaller()
		case ProtobufFormat:
			s.reviewRatingIncrementedMarshaller = s.ReviewRatingIncrementedEventProtobufMarshaller()
		}
	}

	return s.reviewRatingIncrementedMarshaller
}

func (s *serviceLocator) MysqlIntegrationEventOutboxRepository() *mysql.IntegrationReviewEventsOutboxRepository {
	if s.mysqlIntegrationEventOutboxRepository == nil {
		s.mysqlIntegrationEventOutboxRepository = mysql.NewIntegrationReviewEventsOutboxRepository(s.MysqlDB())
	}
	return s.mysqlIntegrationEventOutboxRepository
}

func (s *serviceLocator) MysqlTransactionManager() *mysql.TransactionManager {
	if s.mysqlTransactionManager == nil {
		s.mysqlTransactionManager = mysql.NewTransactionManager(s.MysqlDB())
	}
	return s.mysqlTransactionManager
}

func (s *serviceLocator) MysqlCommandReviewRepository() *mysql.CommandsReviewRepository {
	if s.mysqlCommandsReviewRepository == nil {
		s.mysqlCommandsReviewRepository = mysql.NewReviewRepository(s.MysqlDB())
	}
	return s.mysqlCommandsReviewRepository
}

func (s *serviceLocator) MysqlQueryReviewsRepository() *mysql.QueriesReviewsRepository {
	if s.mysqlQueriesReviewRepository == nil {
		s.mysqlQueriesReviewRepository = mysql.NewGetReviewsRepository(s.MysqlDB())
	}
	return s.mysqlQueriesReviewRepository
}

func (s *serviceLocator) ReviewCreatedEventProducer() *review_created.ReviewCreatedEventProducer {
	if s.reviewCreatedEventProducer == nil {
		s.reviewCreatedEventProducer = review_created.NewReviewCreatedEventProducer(
			s.CreateReviewProducer(),
			s.ReviewCreatedMarshaller(),
		)
	}

	return s.reviewCreatedEventProducer
}

func (s *serviceLocator) ReviewRatingIncrementedEventProducer() *review_rating_incremented.ReviewRatingIncrementedEventProducer {
	if s.reviewRatingIncrementedEventProducer == nil {
		s.reviewRatingIncrementedEventProducer = review_rating_incremented.NewReviewRatingIncrementedEventProducer(
			s.IncrementReviewRatingProducer(),
			s.ReviewRatingIncrementedMarshaller(),
		)
	}

	return s.reviewRatingIncrementedEventProducer
}

func (s *serviceLocator) ReviewCreatedEventConsumer() *review_created.ReviewCreatedEventConsumer {
	if s.reviewCreatedEventConsumer == nil {
		s.reviewCreatedEventConsumer = review_created.NewReviewCreatedEventConsumer(
			s.CreateReviewConsumer(),
			s.ReviewCreatedMarshaller(),
			s.GetReviewsProjector(),
		)
	}

	return s.reviewCreatedEventConsumer
}

func (s *serviceLocator) ReviewRatingIncrementedEventConsumer() *review_rating_incremented.ReviewRatingIncrementedEventConsumer {
	if s.reviewRatingIncrementedEventConsumer == nil {
		s.reviewRatingIncrementedEventConsumer = review_rating_incremented.NewReviewRatingIncrementedEventConsumer(
			s.CreateReviewConsumer(),
			s.ReviewRatingIncrementedMarshaller(),
			s.ReviewRatingIncrementedProjector(),
		)
	}

	return s.reviewRatingIncrementedEventConsumer
}

func (s *serviceLocator) ReviewCreatedKafkaConsumer() *kafka.Consumer {
	if s.reviewCreatedKafkaConsumer == nil {
		var err error
		s.reviewCreatedKafkaConsumer, err = kafka.NewConsumer(s.KafkaURL(), 0, "review_created")
		if err != nil {
			panic(err)
		}
	}
	return s.reviewCreatedKafkaConsumer
}

func (s *serviceLocator) ReviewCreatedKafkaProducer() *kafka.Producer {
	if s.reviewCreatedKafkaProducer == nil {
		var err error
		s.reviewCreatedKafkaProducer, err = kafka.NewProducer(s.KafkaURL(), "review_created")
		if err != nil {
			panic(err)
		}
	}
	return s.reviewCreatedKafkaProducer
}

func (s *serviceLocator) ReviewCreatedPulsarConsumer() *pulsar.Consumer {
	if s.reviewCreatedPulsarConsumer == nil {
		var err error
		s.reviewCreatedPulsarConsumer, err = pulsar.NewConsumer(s.PulsarURL(), "review_created")
		if err != nil {
			panic(err)
		}
	}
	return s.reviewCreatedPulsarConsumer
}

func (s *serviceLocator) ReviewCreatedPulsarProducer() *pulsar.Producer {
	if s.reviewCreatedPulsarProducer == nil {
		var err error
		s.reviewCreatedPulsarProducer, err = pulsar.NewProducer(s.PulsarURL(), "review_created")
		if err != nil {
			panic(err)
		}
	}
	return s.reviewCreatedPulsarProducer
}

func (s *serviceLocator) ReviewCreatedEventJSONMarshaller() *review_created.ReviewCreatedEventJSONMarshaller {
	if s.reviewCreatedEventJSONMarshaller == nil {
		s.reviewCreatedEventJSONMarshaller = &review_created.ReviewCreatedEventJSONMarshaller{}
	}

	return s.reviewCreatedEventJSONMarshaller
}

func (s *serviceLocator) ReviewCreatedEventProtobufMarshaller() *review_created.ReviewCreatedEventProtobufMarshaller {
	if s.reviewCreatedEventProtobufMarshaller == nil {
		s.reviewCreatedEventProtobufMarshaller = &review_created.ReviewCreatedEventProtobufMarshaller{}
	}

	return s.reviewCreatedEventProtobufMarshaller
}

func (s *serviceLocator) ReviewRatingIncrementedEventJSONMarshaller() *review_rating_incremented.ReviewRatingIncrementedEventJSONMarshaller {
	if s.reviewRatingIncrementedEventJSONMarshaller == nil {
		s.reviewRatingIncrementedEventJSONMarshaller = &review_rating_incremented.ReviewRatingIncrementedEventJSONMarshaller{}
	}

	return s.reviewRatingIncrementedEventJSONMarshaller
}

func (s *serviceLocator) ReviewRatingIncrementedEventProtobufMarshaller() *review_rating_incremented.ReviewRatingIncrementedEventProtobufMarshaller {
	if s.reviewRatingIncrementedEventProtobufMarshaller == nil {
		s.reviewRatingIncrementedEventProtobufMarshaller = &review_rating_incremented.ReviewRatingIncrementedEventProtobufMarshaller{}
	}

	return s.reviewRatingIncrementedEventProtobufMarshaller
}

func (s *serviceLocator) ReviewRatingIncrementedKafkaConsumer() *kafka.Consumer {
	if s.reviewRatingIncrementedKafkaConsumer == nil {
		var err error
		s.reviewRatingIncrementedKafkaConsumer, err = kafka.NewConsumer(s.KafkaURL(), 0, "review_rating_incremented")
		if err != nil {
			panic(err)
		}
	}
	return s.reviewRatingIncrementedKafkaConsumer
}

func (s *serviceLocator) ReviewRatingIncrementedKafkaProducer() *kafka.Producer {
	if s.reviewRatingIncrementedKafkaProducer == nil {
		var err error
		s.reviewRatingIncrementedKafkaProducer, err = kafka.NewProducer(s.KafkaURL(), "review_rating_incremented")
		if err != nil {
			panic(err)
		}
	}
	return s.reviewRatingIncrementedKafkaProducer
}

func (s *serviceLocator) ReviewRatingIncrementedPulsarConsumer() *pulsar.Consumer {
	if s.reviewRatingIncrementedPulsarConsumer == nil {
		var err error
		s.reviewRatingIncrementedPulsarConsumer, err = pulsar.NewConsumer(s.PulsarURL(), "review_rating_incremented")
		if err != nil {
			panic(err)
		}
	}
	return s.reviewRatingIncrementedPulsarConsumer
}

func (s *serviceLocator) ReviewRatingIncrementedPulsarProducer() *pulsar.Producer {
	if s.reviewRatingIncrementedPulsarProducer == nil {
		var err error
		s.reviewRatingIncrementedPulsarProducer, err = pulsar.NewProducer(s.PulsarURL(), "review_rating_incremented")
		if err != nil {
			panic(err)
		}
	}
	return s.reviewRatingIncrementedPulsarProducer
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
