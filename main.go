package main

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "cocowork"
	password = "rahasia"
	dbName   = "gophercises_phone"
)

func main() {
	// reset db first
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", psqlInfo)
	// db.Close()
	must(err)
	err = resetDB(db, dbName)
	must(err)

	// connect to new database that has been crated
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbName)
	db, err = sql.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()
	must(db.Ping())
}

func must(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	fmt.Println("Successfully create db")
	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	fmt.Println("Successfully drop db")
	return createDB(db, dbName)
}

// func normalize(phone string) string {
// 	var buf bytes.Buffer

// 	// when iterate, string will be []rune
// 	for _, ch := range phone {
// 		// condition to check if ch >= 0 && ch <= 9
// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch)
// 		}
// 	}

// 	fmt.Println(buf.String())

// 	return buf.String()
// }

func normalize(phone string) string {
	// compile regex to make sure is not error
	// if error it will panic
	// \D -> non digit
	re := regexp.MustCompile("\\D")

	// replace non digit with ""
	return re.ReplaceAllString(phone, "")
}
