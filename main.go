package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/podocarp/mysql-test-test/utils"
)

var testData []utils.Countries = []utils.Countries{}

var db *sql.DB

func main() {
	dsn := fmt.Sprintf("%s:%s@%s/%s?%s",
		"root", "asd", // user, password
		"",     // address, empty for localhost:3306
		"test", // db
		"",     // options
	)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.Exec(`DROP TABLE IF EXISTS countries_bitset;`)
	_, err = db.Exec(`CREATE TABLE countries_bitset (
          id bigint unsigned NOT NULL AUTO_INCREMENT,
          countries BINARY(32),
          PRIMARY KEY (id)
        )`)
	if err != nil {
		panic(err)
	}
	db.Exec(`DROP TABLE IF EXISTS countries_json;`)
	_, err = db.Exec(`CREATE TABLE countries_json (
          id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
          countries JSON,
          PRIMARY KEY (id)
        )`)
	if err != nil {
		panic(err)
	}

	for range 100_000 {
		testData = append(testData, utils.RandomCountries())
	}

	fmt.Println("Starting test now")

	timeIt(WriteAsBitsets, "Write using bitset")
	timeIt(WriteAsJSON, "Write using json")
	// timeIt(ReadAsBitsets, "Read using bitset")
	// timeIt(ReadAsJSON, "Read using json")

	// profileIt(WriteAsJSON, "json.cpu")
	// profileIt(WriteAsBitsets, "bitset.cpu")
}

func timeIt(fun func(), msg string) {
	now := time.Now()
	fun()
	t := time.Since(now)
	fmt.Printf("%s time taken: %v\n", msg, t)
}

func profileIt(fun func(), filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}

	fun()

	pprof.StopCPUProfile()
	fmt.Println("profile", filename, "done")
}

func WriteAsBitsets() {
	for _, data := range testData {
		_, err := db.Exec(`INSERT INTO countries_bitset (countries) VALUES (?)`, &data)
		if err != nil {
			panic(err)
		}
	}
}

type Row struct {
	ID        uint64
	Countries *utils.Countries
}

func ReadAsBitsets() {
	rows, err := db.Query(`SELECT id,countries FROM countries_bitset`)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var row Row
		rows.Scan(&row.ID, &row.Countries)
	}
}

func WriteAsJSON() {
	for _, data := range testData {
		j, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		db.Exec(`INSERT INTO countries_json (countries) VALUES (?)`, string(j))
	}
}

func ReadAsJSON() {
	rows, err := db.Query(`SELECT id,countries FROM countries_json`)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var row Row
		var j []byte
		rows.Scan(&row.ID, &j)
		err := json.Unmarshal(j, &row.Countries)
		if err != nil {
			panic(err)
		}
	}
}
