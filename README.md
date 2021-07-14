# whosonfirst2pgsql

Convert sqlite data to postgresql

### MAIN OPTIONS
#### -a, --append
Run in append mode. Adds the OSM change file into the database without removing existing data.
#### -c, --create
Run in create mode. This is the default if **-a**, **--append** is not specified. Removes existing data from the database tables!

### LOGGING OPTIONS
#### --log-level=LEVEL
Set log level (‘debug’, ‘info’ (default), ‘warn’, or ‘error’).
#### --log-progress=VALUE
Enable (`true`) or disable (`false`) progress logging. Setting this to `auto` will enable progress logging on the console and disable it if the output is redirected to a file. Default: true.
#### --log-sql
Enable logging of SQL commands for debugging.
#### --log-sql-data
Enable logging of all data added to the database. This will write out a huge amount of data! For debugging.
#### -v, --verbose
Same as `--log-level=debug`.


### DATABASE OPTIONS
#### -d, --database=NAME
The name of the PostgreSQL database to connect to. If this parameter contains an `=` sign or starts with a valid URI prefix (`postgresql://` or `postgres://`), it is treated as a conninfo string. See the PostgreSQL manual for details.
#### -U, --username=NAME
Postgresql user name.
#### -W, --password
Force password prompt.
#### -H, --host=HOSTNAME
Database server hostname or unix domain socket location.
#### -P, --port=PORT
Database server port.


### INPUT OPTIONS
#### -r, --input-reader=FORMAT
Select format of the input file. Available choices are auto (default) for autodetecting the format, xml for OSM XML format files, o5m for o5m formatted files and pbf for OSM PBF binary format.
#### -b, --bbox=MINLON,MINLAT,MAXLON,MAXLAT
Apply a bounding box filter on the imported data. Example: --bbox -0.5,51.25,0.5,51.75

### PGSQL OUTPUT OPTIONS
#### -i, --tablespace-index=TABLESPC
Store all indexes in the PostgreSQL tablespace TABLESPC. This option also affects the middle tables.
#### --tablespace-main-data=TABLESPC
Store the data tables in the PostgreSQL tablespace TABLESPC.
#### --tablespace-main-index=TABLESPC
Store the indexes in the PostgreSQL tablespace TABLESPC.
#### --latlong
Store coordinates in degrees of latitude & longitude.
#### -m, --merc
Store coordinates in Spherical Mercator (Web Mercator, EPSG:3857) (the default).
#### -E, --proj=SRID
Use projection EPSG:SRID.
#### -p, --prefix=PREFIX
Prefix for table names (default: planet_osm). This option affects the middle as well as the pgsql output table names.
#### -k, --hstore
Add tags without column to an additional hstore (key/value) column in the database tables.
#### -j, --hstore-all
Add all tags to an additional hstore (key/value) column in the database tables.
#### -z, --hstore-column=PREFIX
Add an additional hstore (key/value) column named PREFIX containing all tags that have a key starting with PREFIX, eg \--hstore-column "name:" will produce an extra hstore column that contains all name:xx tags.
#### --hstore-match-only
Only keep objects that have a value in at least one of the non-hstore columns.
#### --hstore-add-index
Create indexes for all hstore columns after import.
#### --output-pgsql-schema=SCHEMA
Use PostgreSQL schema SCHEMA for all tables, indexes, and functions in the pgsql output (default is no schema, i.e. the public schema is used).