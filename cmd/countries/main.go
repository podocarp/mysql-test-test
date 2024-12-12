package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/podocarp/mysql-test-test/db"
	"github.com/podocarp/mysql-test-test/utils"
)

var testData []utils.Countries = []utils.Countries{}

const (
	VERBOSE = false
	ITERS   = 30
)

var pool *sql.DB

func main() {
	TestWrite()
	TestRead()
}

func TestWrite() {
	jtimer := utils.NewTimer("json").SetSilent()
	btimer := utils.NewTimer("bitset").SetSilent()
	for range ITERS {
		InitTest(1000)
		jtimer.TimeIt(WriteAsJSON)
	}
	for range ITERS {
		InitTest(1000)
		btimer.TimeIt(WriteAsBitsets)
	}
	jtimer.Echo()
	btimer.Echo()
	utils.GraphTimers("countries-w.html", "JSON vs Bitset (Writing)", jtimer, btimer)
}

func TestRead() {
	jtimer := utils.NewTimer("json").SetSilent()
	btimer := utils.NewTimer("bitset").SetSilent()
	for range ITERS {
		InitTest(1000)
		jtimer.TimeIt(ReadAsJSON)
	}
	for range ITERS {
		InitTest(1000)
		btimer.TimeIt(ReadAsBitsets)
	}
	jtimer.Echo()
	btimer.Echo()
	utils.GraphTimers("countries-r.html", "JSON vs Bitset (Reading)", jtimer, btimer)
}

func ClearTables() {
	var err error
	_, err = pool.Exec(`CREATE TABLE IF NOT EXISTS countries_bitset (
          id bigint unsigned NOT NULL AUTO_INCREMENT,
          countries BINARY(32),
          PRIMARY KEY (id)
        ) AUTO_INCREMENT=0 ENGINE=InnoDB`)
	if err != nil {
		panic(err)
	}
	pool.Exec(`TRUNCATE TABLE countries_bitset;`)

	_, err = pool.Exec(`CREATE TABLE IF NOT EXISTS countries_json (
          id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
          countries JSON,
          PRIMARY KEY (id)
        ) AUTO_INCREMENT=0 ENGINE=InnoDB`)
	if err != nil {
		panic(err)
	}
	pool.Exec(`TRUNCATE TABLE countries_json;`)
	fmt.Println("Cleared all tables.")
}

func InitTest(numentries int) {
	testData = make([]utils.Countries, numentries)
	for i := range numentries {
		testData[i] = utils.RandomCountries()
	}
}

type Row struct {
	ID        uint64
	Countries *utils.Countries
}

func WriteAsBitsets() {
	for i, data := range testData {
		_, err := pool.Exec(`INSERT INTO countries_bitset (countries) VALUES (?);`, &data)
		if err != nil {
			panic(err)
		}
		if VERBOSE && i > 0 && i%10000 == 0 {
			fmt.Println("\tWritten", i, "rows")
		}
	}
}

func ReadAsBitsets() {
	for i := range len(testData) {
		row := pool.QueryRow(`SELECT id,countries FROM countries_bitset WHERE id=?;`, i+1)
		var r Row
		err := row.Scan(&r.ID, &r.Countries)
		if err != nil {
			panic(err)
		}
	}
}

func WriteAsJSON() {
	for i, data := range testData {
		j, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		_, err = pool.Exec(`INSERT INTO countries_json (countries) VALUES (?);`, string(j))
		if err != nil {
			panic(err)
		}
		if VERBOSE && i > 0 && i%10000 == 0 {
			fmt.Println("\tWritten", i, "rows")
		}
	}
}

func ReadAsJSON() {
	for i := range len(testData) {
		row := pool.QueryRow(`SELECT id,countries FROM countries_json WHERE id=?;`, i+1)
		var r Row
		var j []byte
		err := row.Scan(&r.ID, &j)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(j, &r.Countries)
		if err != nil {
			panic(i)
		}
	}
}

func init() {
	pool = db.Connect()
	ClearTables()
}
