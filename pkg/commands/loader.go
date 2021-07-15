package commands

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

type PasswordProvider interface {
	ReadPassword() (string, error)
}
type Config struct {
	File                       string
	Create                     bool
	DbName                     string
	UserName                   string
	Password                   string
	Host                       string
	Port                       int
	UseLatLng                  bool
	UseGeom                    bool
	TableName                  string
	InclKeyValues              bool
	ExcludeColumnFromKeyValues bool
	UseJson                    bool
	MatchOnly                  bool
	Schema                     string
	WorkerCount                int
}
type Loader struct {
	PasswordProvider PasswordProvider
	Config           Config
	Runner           Runner
}

type TerminalPasswordReader struct{}

func (pr TerminalPasswordReader) ReadPassword() (string, error) {
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	return string(password), err
}

func (t *Loader) DoLoad(context *cli.Context) error {
	t.Config = Config{}
	if context.NArg() >= 1 {
		t.Config.File = context.Args().Get(0)
	}
	t.Config.Create = context.Bool("c")
	if context.Bool("a") {
		t.Config.Create = false
	}
	t.Config.DbName = context.String("d")
	t.Config.UserName = context.String("U")
	if context.Bool("W") {
		// should prompt password and set it
		log.Println("Now, please type in the password (mandatory): ")
		pwd, _ := t.PasswordProvider.ReadPassword()
		t.Config.Password = pwd
	} else {
		t.Config.Password = "mysecretpassword"
	}
	t.Config.Host = context.String("H")
	t.Config.Port = context.Int("P")
	t.Config.UseGeom = true
	if context.Bool("latlong") {
		t.Config.UseLatLng = true
		t.Config.UseGeom = false
	}
	t.Config.TableName = context.String("p")
	t.Config.InclKeyValues =
		context.Bool("hstore") ||
			context.Bool("hstore-all") ||
			context.Bool("json") ||
			context.Bool("json-all")

	t.Config.ExcludeColumnFromKeyValues = context.Bool("hstore") || context.Bool("json")
	t.Config.UseJson = context.Bool("json") || context.Bool("json-all")
	t.Config.MatchOnly = context.Bool("match-only")
	t.Config.Schema = context.String("schema")
	t.Config.TableName = context.String("t")
	t.Config.WorkerCount = context.Int("workers")
	return t.Runner.Run(&t.Config)
}
