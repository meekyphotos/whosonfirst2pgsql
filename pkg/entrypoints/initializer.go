package entrypoints

import (
	"github.com/meekyphotos/whosonfirst2pgsql/pkg/commands"
	"github.com/urfave/cli/v2"
)

func InitApp(loader *commands.Loader) *cli.App {
	return &cli.App{
		Name:        "whosonfirst2pgsql",
		Description: "Convert data to pgsql",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "a", Aliases: []string{"append"}, Value: false, Usage: "Run in append mode. Adds the OSM change file into the database without removing existing data."},
			&cli.BoolFlag{Name: "c", Aliases: []string{"create"}, Value: true, Usage: "Run in create mode. This is the default if -a, --append is not specified. Removes existing data from the database tables!"},
			// DATABASE OPTIONS
			&cli.StringFlag{Name: "d", Aliases: []string{"database"}, Value: "", Required: true, Usage: "The name of the PostgreSQL database to connect to. If this parameter contains an = sign or starts with a valid URI prefix (postgresql:// or postgres://), it is treated as a conninfo string. See the PostgreSQL manual for details."},
			&cli.StringFlag{Name: "U", Aliases: []string{"username"}, Value: "postgres", Usage: "Postgresql user name."},
			&cli.BoolFlag{Name: "W", Aliases: []string{"password"}, Value: false, Usage: "Force password prompt."},
			&cli.StringFlag{Name: "H", Aliases: []string{"host"}, Value: "localhost", Usage: "Database server hostname or unix domain socket location."},
			&cli.IntFlag{Name: "P", Aliases: []string{"port"}, Value: 5432, Usage: "Database server port."},
			&cli.IntFlag{Name: "workers", Value: 4, Usage: "Number of workers"},

			// OUTPUT FORMAT
			&cli.BoolFlag{Name: "latlong", Value: false, Usage: "Store coordinates in degrees of latitude & longitude."},
			&cli.StringFlag{Name: "t", Aliases: []string{"table"}, Value: "planet_data", Usage: "Output table name"},

			&cli.BoolFlag{Name: "k", Aliases: []string{"hstore"}, Value: false, Usage: "Add tags without column to an additional hstore (key/value) column in the database tables."},
			&cli.BoolFlag{Name: "j", Aliases: []string{"hstore-all"}, Value: false, Usage: "Add all tags to an additional hstore (key/value) column in the database tables."},
			&cli.BoolFlag{Name: "json", Value: false, Usage: "Add tags without column to an additional json (key/value) column in the database tables."},
			&cli.BoolFlag{Name: "json-all", Value: false, Usage: "Add all tags to an additional json (key/value) column in the database tables."},

			&cli.BoolFlag{Name: "match-only", Value: false, Usage: "Only keep objects that have a value in at least one of the non-key value columns."},

			&cli.StringFlag{Name: "schema", Value: "public", Usage: "Use PostgreSQL schema SCHEMA for all tables, indexes, and functions in the pgsql output (default is no schema, i.e. the public schema is used)."},
		},
		UseShortOptionHandling: true,
		Action:                 loader.DoLoad,
	}
}
