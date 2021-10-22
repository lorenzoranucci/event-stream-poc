package main

import (
	"fmt"
	"os"

	"github.com/ProntoPro/golang-kafka/internal/cmd"
)

var (
	version = "dev"
	app     = cmd.GetApp(version)
)

func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}
