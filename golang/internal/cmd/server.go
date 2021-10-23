package cmd

import (
	"github.com/urfave/cli"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/create_review"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/in_memory"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/kafka"
	kafka2 "github.com/ProntoPro/event-stream-golang/pkg/kafka"
)

func getServerCommand(baseFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:   "server",
		Action: runServer,
		Usage:  "Run the http server which expose the POC API",
		Flags: append(
			baseFlags,
			cli.IntFlag{
				Name:   "port",
				Value:  8080,
				Usage:  "Server port",
				EnvVar: "PPRO_PORT",
			},
		),
	}
}

func runServer(c *cli.Context) error {
	kafkaProducer, err := kafka2.NewProducer(c.String("kafka-url"))
	if err != nil {
		return err
	}

	server := http.NewServer(
		c.Int("port"),
		create_review.NewCreateReviewHandler(
			application.NewCreateReviewCommandHandler(
				&in_memory.ReviewRepository{},
				kafka.NewReviewCreatedEventBus(
					kafkaProducer,
				),
			),
		),
	)

	server.Run()

	return nil
}
