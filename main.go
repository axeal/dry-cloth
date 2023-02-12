package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		fmt.Println("Not yet implemented")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
