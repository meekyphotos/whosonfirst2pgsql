package commands

import (
	"github.com/urfave/cli/v2"
	"log"
)

type Loader struct {
}

func (t *Loader) DoLoad(context *cli.Context) error {
	if context.NArg() >= 1 {
		log.Println(context.Args().Get(0))
	}
	return nil
}
