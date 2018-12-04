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
	must(err)
	err = resetDB(db, dbName)
	must(err)
	db.Close()

	// connect to new database that has been crated
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbName)
	db, err = sql.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()
	must(db.Ping())
	must(createPhoneNumbersTable(db))

	id, err := insertPhoneNumber(db, "123456789")
	must(err)
	id, err = insertPhoneNumber(db, "123 456 7891")
	must(err)
	id, err = insertPhoneNumber(db, "(123) 456 7892")
	must(err)
	id, err = insertPhoneNumber(db, "(123) 456-7893")
	must(err)
	id, err = insertPhoneNumber(db, "123-456-7894")
	must(err)
	id, err = insertPhoneNumber(db, "123-456-7890")
	must(err)
	id, err = insertPhoneNumber(db, "1234567892")
	must(err)
	id, err = insertPhoneNumber(db, "(123) 456-7890")
	must(err)
	fmt.Println("id : ", id)
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

func createPhoneNumbersTable(db *sql.DB) error {
	statement := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			phone_number VARCHAR(255)
		)`)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	fmt.Println("Successfully create phone_numbers table")
	return nil
}

func insertPhoneNumber(db *sql.DB, phone string) (int, error) {
	var id int
	statement := `INSERT INTO phone_numbers (phone_number) VALUES($1) RETURNING id`
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
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
