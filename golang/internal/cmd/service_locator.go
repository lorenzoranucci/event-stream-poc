package cmd

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/domain"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/patch_review"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/commands"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/queries"
	get_reviews2 "github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/projectors/get_reviews"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/create_review"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/get_reviews"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/mysql"
	"github.com/ProntoPro/event-stream-golang/pkg/kafka"
)

func newServiceLocator(
	kafkaURL string,
	mysqlURL string,
) *serviceLocator {
	return &serviceLocator{
		kafkaURL: kafkaURL,
		mysqlURL: mysqlURL,
	}
}

type serviceLocator struct {
	createReviewCommandHandler          *commands.CreateReviewCommandHandler
	getReviewsQueryHandler              *queries.GetReviewsQueryHandler
	incrementReviewRatingCommandHandler *commands.IncrementReviewRatingCommandHandler

	createReviewHandler *create_review.CreateReviewHandler
	getReviewsHandler   *get_reviews.GetReviewsHandler
	patchReviewHandler  *patch_review.PatchReviewHandler

	mysqlCommandsReviewRepository *mysql.CommandsReviewRepository
	mysqlQueriesReviewRepository  *mysql.QueriesReviewsRepository

	reviewCreatedProjector           *get_reviews2.ReviewCreatedProjector
	reviewRatingIncrementedProjector *get_reviews2.ReviewRatingIncrementedProjector

	reviewCreatedKafkaConsumer *kafka.Consumer
	reviewCreatedKafkaProducer *kafka.Producer

	reviewRatingIncrementedKafkaConsumer *kafka.Consumer
	reviewRatingIncrementedKafkaProducer *kafka.Producer

	mysqlDB *sqlx.DB

	kafkaURL        string
	mysqlURL        string
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
		)
	}

	return s.incrementReviewRatingCommandHandler
}

func (s *serviceLocator) CommandReviewRepository() domain.ReviewRepository {
	return s.MysqlCommandReviewRepository()
}

func (s *serviceLocator) QueryReviewRepository() queries.ReviewRepository {
	return s.MysqlQueryReviewsRepository()
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
