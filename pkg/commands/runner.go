package commands

import (
	"archive/tar"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/godruoyi/go-snowflake"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/valyala/fastjson"
)

type Runner interface {
	Run(c *Config) error
}

type CmdRunner struct{}

func (c *CmdRunner) Run(config *Config) error {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
		config.UserName, config.DbName, config.Password, config.Host)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	db := DbConn{
		cfg:  config,
		conn: conn,
		written: progressbar.Default(
			-1,
			"Writing",
		),
	}

	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)
	err = db.InitTable()
	if err != nil {
		return err
	}

	log.Println("Start reading... ", config.File)
	contentChannel := make(chan []byte, 10000)
	reqChannel := make(chan Req, 10000)
	var localError error
	go func() {
		err := scanTarFile(contentChannel, &db, config)
		if err != nil {
			localError = err
		}
	}()

	pool := fastjson.ParserPool{}

	var jsonWorkers sync.WaitGroup
	var pgWorkers sync.WaitGroup

	for i := 0; i < config.WorkerCount; i++ {
		jsonWorkers.Add(1)
		go func(index int) {
			err := jsonParser(contentChannel, reqChannel, &pool)
			if err != nil {
				localError = err
			}
			jsonWorkers.Done()
		}(i)
	}
	batchedRequests := batchRequest(reqChannel, 10000, time.Second)

	pgWorkers.Add(1)
	go func() {
		err := dbLoader(batchedRequests, &db)
		if err != nil {
			localError = err
		}
		pgWorkers.Done()
	}()

	jsonWorkers.Wait()
	close(reqChannel)
	pgWorkers.Wait()
	if localError == nil {
		log.Println("Creating index")
		err := db.CreateIndexes()
		if err != nil {
			return err
		}
	}
	return localError
}

func batchRequest(values <-chan Req, maxItems int, maxTimeout time.Duration) chan []Req {
	batches := make(chan []Req)

	go func() {
		defer close(batches)

		for keepGoing := true; keepGoing; {
			var batch []Req
			expire := time.After(maxTimeout)
			for {
				select {
				case value, ok := <-values:
					if !ok {
						keepGoing = false
						goto done
					}

					batch = append(batch, value)
					if len(batch) == maxItems {
						goto done
					}

				case <-expire:
					goto done
				}
			}

		done:
			if len(batch) > 0 {
				batches <- batch
			}
		}
	}()

	return batches
}

func scanTarFile(channel chan []byte, db *DbConn, config *Config) error {
	open, err := os.Open(config.File)
	if err != nil {
		return err
	}

	defer func(open *os.File) {
		err := open.Close()
		if err != nil {
			panic(err)
		}
	}(open)

	reader := tar.NewReader(open)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if header.Typeflag == tar.TypeReg && strings.HasSuffix(header.Name, "geojson") {
			content := make([]byte, header.Size)
			_, err := reader.Read(content)
			if err != nil && err != io.EOF {
				close(channel)
				return err
			}

			channel <- content
			db.read++
			if db.read%1000 == 0 {
				db.written.ChangeMax64(db.read)
			}
		}
	}
	db.written.ChangeMax64(db.read)
	close(channel)
	return nil
}

func jsonParser(channel chan []byte, reqChannel chan Req, pool *fastjson.ParserPool) error {

	for {
		select {
		case content := <-channel:
			if len(content) == 0 {
				return nil
			}
			parser := pool.Get()
			p, err := parser.ParseBytes(content)
			if err != nil {
				return err
			}
			req := fromJsonValue(p)
			reqChannel <- req
			pool.Put(parser)
		default:
		}
	}
}

func dbLoader(channel chan []Req, db *DbConn) error {
	for {
		select {
		case content := <-channel:
			if len(content) == 0 {
				return nil
			}
			err := db.ProcessRequest(content)
			if err != nil {
				return err
			}
		default:
		}
	}
}

type Req struct {
	id             int64
	latitude       float64
	longitude      float64
	countryCode    string
	preferredNames map[string]string
	variantNames   map[string]string
	macroregionId  int64
	continentId    int64
	countryId      int64
	countyId       int64
	localityId     int64
	regionId       int64
	metadata       map[string]string
}

type DbConn struct {
	conn         *sql.DB
	cfg          *Config
	tableColumns []string
	read         int64
	written      *progressbar.ProgressBar
}

var fieldType = map[string]string{
	"id":              "bigint",
	"wof_id":          "bigint",
	"latitude":        "double precision",
	"longitude":       "double precision",
	"continent_id":    "bigint",
	"country_id":      "bigint",
	"country_code":    "text",
	"county_id":       "bigint",
	"locality_id":     "bigint",
	"region_id":       "bigint",
	"preferred_names": "jsonb",
	"variant_names":   "jsonb",
	"metadata":        "jsonb",
	"latlng":          "geography(POINT)",
}
var fields = []string{
	"id",
	"wof_id",
	"continent_id",
	"country_id",
	"country_code",
	"county_id",
	"locality_id",
	"region_id",
	"preferred_names",
	"variant_names",
}

var latLngFields = []string{
	"latitude",
	"longitude",
}

var geomFields = []string{
	"latlng",
}

func (db *DbConn) InitTable() error {
	// save for future reference
	txn, err := db.conn.Begin()
	if err != nil {
		return err
	}

	if db.cfg.Create {
		_, err := txn.Exec("DROP TABLE IF EXISTS " + db.cfg.Schema + "." + db.cfg.TableName)
		if err != nil {
			return err
		}
	}
	db.tableColumns = make([]string, 0)
	db.tableColumns = append(db.tableColumns, fields...)
	if db.cfg.InclKeyValues {
		db.tableColumns = append(db.tableColumns, "metadata")
	}
	if db.cfg.UseGeom {
		db.tableColumns = append(db.tableColumns, geomFields...)
	} else {
		db.tableColumns = append(db.tableColumns, latLngFields...)
	}

	stmt := strings.Builder{}
	stmt.WriteString("CREATE TABLE IF NOT EXISTS ")
	stmt.WriteString(db.cfg.Schema)
	stmt.WriteString(".")
	stmt.WriteString(db.cfg.TableName)
	stmt.WriteString("\n(")
	for i, f := range db.tableColumns {
		stmt.WriteString(f)
		stmt.WriteString(" ")
		stmt.WriteString(fieldType[f])
		if f == "id" {
			stmt.WriteString(" primary key")
		}
		if i < len(db.tableColumns)-1 {
			stmt.WriteString(", \n")
		}
	}
	stmt.WriteString("\n)")
	sql := stmt.String()
	_, err = txn.Exec(sql)
	if err != nil {
		return txn.Rollback()
	}
	return txn.Commit()
}
func (db *DbConn) CreateIndexes() error {

	if db.cfg.UseGeom {
		txn, err := db.conn.Begin()
		if err != nil {
			return err
		}
		txn.Exec("CREATE INDEX ON " + db.cfg.Schema + "." + db.cfg.TableName + " using BRIN (latlng)")
		txn.Commit()
	}
	return nil
}

func (db *DbConn) ProcessRequest(reqs []Req) error {
	txn, err := db.conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := txn.Prepare(pq.CopyIn(db.cfg.TableName,
		db.tableColumns...,
	))
	if err != nil {
		return err
	}
	for _, r := range reqs {
		vals := make([]interface{}, len(db.tableColumns))
		for i, f := range db.tableColumns {
			switch f {
			case "id":
				vals[i] = snowflake.ID()
			case "wof_id":
				vals[i] = r.id
			case "latitude":
				vals[i] = r.latitude
			case "longitude":
				vals[i] = r.longitude
			case "continent_id":
				vals[i] = r.continentId
			case "country_id":
				vals[i] = r.countryId
			case "country_code":
				vals[i] = r.countryCode
			case "county_id":
				vals[i] = r.countyId
			case "locality_id":
				vals[i] = r.localityId
			case "region_id":
				vals[i] = r.regionId
			case "preferred_names":
				prefNames, _ := json.Marshal(r.preferredNames)
				vals[i] = string(prefNames)
			case "variant_names":
				variantNames, _ := json.Marshal(r.variantNames)
				vals[i] = string(variantNames)
			case "metadata":
				metadata, _ := json.Marshal(r.metadata)
				vals[i] = string(metadata)
			case "latlng":
				if r.latitude != 0 && r.longitude != 0 {
					vals[i] = fmt.Sprintf("SRID=4326;POINT(%f %f)", r.latitude, r.longitude)
				}
			}
		}
		stmt.Exec(vals...)
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	_ = db.written.Add(len(reqs))
	return txn.Commit()
}
