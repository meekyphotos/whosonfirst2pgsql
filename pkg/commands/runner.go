package commands

import (
	"archive/tar"
	"database/sql"
	"github.com/lib/pq"
	"github.com/valyala/fastjson"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type Runner interface {
	Run(c *Config) error
}

type CmdRunner struct{}

func (c *CmdRunner) Run(config *Config) error {
	log.Println("Start reading... ", config.File)
	pool := fastjson.ParserPool{}
	contentChannel := make(chan []byte, 10000)
	var wg sync.WaitGroup
	var localError error
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(index int) {
			log.Println("Start Worker #", index, "...")
			err := jsonParser(contentChannel, &pool)
			if err != nil {
				localError = err
			}
			wg.Done()
			log.Println("Worker #", index, " job completed")
		}(i)

	}

	go func() {
		err := scanTarFile(contentChannel, config)
		if err != nil {
			localError = err
		}
	}()

	wg.Wait()
	return localError
}

func scanTarFile(channel chan []byte, config *Config) error {
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
		} else {
		}
	}
	close(channel)
	return nil
}

func jsonParser(channel chan []byte, pool *fastjson.ParserPool) error {

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
			// do parsing logic

			p.GetObject("properties").Visit(func(key []byte, v *fastjson.Value) {
				keyString := string(key)
				if v.Type() == fastjson.TypeNumber {
					// parse numbers into strings
					// map keyString
					val := v.GetFloat64()

				} else if v.Type() == fastjson.TypeString {
					// map keyString
					val := string(v.GetStringBytes())
				}
			})
			//log.Println(name)
			pool.Put(parser)
		default:
			break

		}

	}
}

type DbConn struct {
	conn *sql.DB
	cfg  *Config
}
type Req struct {
	id              int64
	latitude        float64
	longitude       float64
	areaSquareM     float64
	countryCode     string
	preferredNames  map[string]string
	variantNames    map[string]string
	continentId     int64
	countryId       int64
	countyId        int64
	localityId      int64
	regionId        int64
	otherAttributes map[string]string
}

func (db *DbConn) ProcessRequest(reqs []Req) error {
	txn, err := db.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn(db.cfg.TableName,
		"id",
		"latitude",
		"longitude",
		"area_square_m",
		"country_code",
		"preferred_names",
		"variant_names",
		"continent_id",
		"country_id",
		"county_id",
		"locality_id",
		"region_id",
	))

}
