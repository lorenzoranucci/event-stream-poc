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
			EnvVar: "PPRO_KAFKA_URL",
			Usage:  "Kafka url",
		},
	}

	app.Commands = []cli.Command{
		getServerCommand(app.Flags),
	}

	return app
}
