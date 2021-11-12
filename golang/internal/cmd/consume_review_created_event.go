package cmd

import (
	"github.com/urfave/cli"
)

func consumeReviewCreatedEvent(baseFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:   "consume-review-created",
		Action: consumeReviewCreated,
		Flags:  baseFlags,
	}
}

func consumeReviewCreated(c *cli.Context) error {
	/*serviceLocator := newServiceLocatorFromCliContext(c)

	migrate := Migrate{db: serviceLocator.MysqlDB()}
	migrate.Migrate()

	consumer := serviceLocator.ReviewCreatedEventConsumer()

	err := consumer.Consume()
	if err != nil {
		fmt.Println(err)
	}*/

	return nil
}
