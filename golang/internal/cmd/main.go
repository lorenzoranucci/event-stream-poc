package cmd

import (
	"github.com/urfave/cli"
)

func GetApp(version string) *cli.App {
	app := cli.NewApp()

	app.Version = version

	app.Name = "ProntoPro event-stream-golang poc"
	app.Usage = ""

	app.HideVersion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "kafka-url",
			EnvVar: "KAFKA_URL",
		},
		cli.StringFlag{
			Name:   "messaging-protocol",
			EnvVar: "MESSAGING_PROTOCOL",
		},
	}

	app.Commands = []cli.Command{
		getServerCommand(app.Flags),
		consumeReviewCreatedEvent(app.Flags),
	}

	return app
}
