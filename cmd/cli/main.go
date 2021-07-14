package main

import (
	"github.com/meekyphotos/whosonfirst2pgsql/pkg/commands"
	"github.com/meekyphotos/whosonfirst2pgsql/pkg/entrypoints"
	"log"
	"os"
)

func main() {
	loader := commands.Loader{
		PasswordProvider: commands.TerminalPasswordReader{},
		Runner:           &commands.CmdRunner{},
	}

	app := entrypoints.InitApp(&loader)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
