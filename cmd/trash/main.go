package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/podocarp/mysql-test-test/db"
	"github.com/podocarp/mysql-test-test/utils"
)

var pool *sql.DB

func main() {
	fmt.Println("Starting test: writing junk")
	timer := utils.NewTimer("write")
	for range 30 {
		timer.TimeIt(func() { WriteTrash(1000) })
	}
	utils.GraphTimers("write-junk.html", "Writing Junk", timer)
}

func RandomString() string {
	nBig, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		panic(err)
	}
	n := nBig.Int64()
	b := make([]byte, n)
	_, err = rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func WriteTrash(iterations int) {
	for range iterations {
		trash := RandomString()
		_, err := pool.Exec(`INSERT INTO junk_test (trash) VALUES (?);`, trash)
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	pool = db.Connect()
	_, err := pool.Exec(`CREATE TABLE IF NOT EXISTS junk_test (
          id bigint unsigned NOT NULL AUTO_INCREMENT,
          trash JSON,
          PRIMARY KEY (id)
        ) AUTO_INCREMENT=0 ENGINE=InnoDB`)
	if err != nil {
		panic(err)
	}
	_, err = pool.Exec(`TRUNCATE TABLE junk_test`)
	if err != nil {
		panic(err)
	}
}
