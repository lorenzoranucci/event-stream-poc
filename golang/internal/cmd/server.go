package cmd

import (
	"github.com/ProntoPro/golang-kafka/internal/pkg/application"
	"github.com/ProntoPro/golang-kafka/internal/pkg/infrastructure/http"
	"github.com/ProntoPro/golang-kafka/internal/pkg/infrastructure/http/create_review"
	"github.com/ProntoPro/golang-kafka/internal/pkg/infrastructure/in_memory"
	"github.com/ProntoPro/golang-kafka/internal/pkg/infrastructure/kafka"
	"github.com/urfave/cli"
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
	server := http.NewServer(
		c.Int("port"),
		create_review.NewCreateReviewHandler(
			application.NewCreateReviewCommandHandler(
				&in_memory.ReviewRepository{},
				&kafka.ReviewCreatedEventBus{},
			),
		),
	)

	server.Run()

	return nil
}
