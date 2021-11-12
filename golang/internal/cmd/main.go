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
			Name:   "event_stream-url",
			EnvVar: "KAFKA_URL",
		},
		cli.StringFlag{
			Name:   "mysql-url",
			EnvVar: "MYSQL_URL",
		},
	}

	app.Commands = []cli.Command{
		getServerCommand(app.Flags),
		consumeReviewCreatedEvent(app.Flags),
	}

	return app
}

func newServiceLocatorFromCliContext(c *cli.Context) *serviceLocator {
	return newServiceLocator(
		c.String("kafka-url"),
		c.String("mysql-url"),
	)
}
