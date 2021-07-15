# whosonfirst2pgsql

### NAME:
whosonfirst2pgsql - A new cli application

### USAGE:
whosonfirst2pgsql [global options] command [command options] [arguments...]

#### DESCRIPTION:
Convert data to pgsql

### COMMANDS:
help, h  Shows a list of commands or help for one command

### GLOBAL OPTIONS:
|Flag|Description|
|---|---|
|-a, --append                |Run in append mode. Adds the OSM change file into the database without removing existing data. (default: false)|
|-c, --create                |Run in create mode. This is the default if -a, --append is not specified. Removes existing data from the database tables! (default: true)|
|-d value, --database `value`  |The name of the PostgreSQL database to connect to. If this parameter contains an = sign or starts with a valid URI prefix (postgresql:// or postgres://), it is t reated as a conninfo string. See the PostgreSQL manual for details.|
|-U value, --username `value`  |Postgresql user name. (default: "postgres")|
|-W, --password              |Force password prompt. (default: false)|
|-H value, --host `value`      |Database server hostname or unix domain socket location. (default: "localhost")|
|-P value, --port `value`     |Database server port. (default: 5432)|
|--workers `value`             |Number of workers (default: 4)|
|--latlong                   |Store coordinates in degrees of latitude & longitude. (default: false)|
|-t value, --table `value`     |Output table name (default: "planet_data")|
|-j, --json                  |Add tags without column to an additional json (key/value) column in the database tables. (default: false)|
|--schema value              |Use PostgreSQL schema SCHEMA for all tables, indexes, and functions in the pgsql output (default is no schema, i.e. the public schema is used). (default: "public")|
|--help, -h                  |show help (default: false)|

