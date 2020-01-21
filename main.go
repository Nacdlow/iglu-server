package main

import (
	"gitlab.com/group-nacdlow/nacdlow-server/cmd"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// VERSION specifies the version of nacdlow-server
var VERSION = "0.1.0"

func main() {
	app := cli.NewApp()
	app.Name = "nacdlow-server"
	app.Usage = "The smart home system server application"
	app.Version = VERSION
	app.Commands = []*cli.Command{
		cmd.CmdStart,
		cmd.CmdAdduser,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
