package cmd

import (
	"github.com/urfave/cli"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http"
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
	serviceLocator, err := NewServiceLocator(c.String("kafka-url"))
	if err != nil {
		return err
	}

	useJSON := true
	switch c.String("messaging-protocol") {
	case "protobuf":
		useJSON = false
	}

	handler := serviceLocator.CreateReviewHandlerWithKafkaAndJSON()
	if !useJSON {
		handler = serviceLocator.CreateReviewHandlerWithKafkaAndProtobuf()
	}

	server := http.NewServer(
		c.Int("port"),
		handler,
	)

	server.Run()

	return nil
}
