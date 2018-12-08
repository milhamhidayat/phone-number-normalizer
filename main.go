package main

import (
	"database/sql"
	"fmt"
	db "phone-number-normalizer/db"
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

type phone struct {
	ID     int
	Number string
}

func main() {
	// reset db first
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	must(db.Reset("postgres", psqlInfo, dbName))

	// connect to new database that has been crated
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbName)
	must(db.Migrate("postgres", psqlInfo))

	_, err = insertPhoneNumber(db, "123456789")
	must(err)
	_, err = insertPhoneNumber(db, "123 456 7891")
	must(err)
	id, err := insertPhoneNumber(db, "(123) 456 7892")
	must(err)
	_, err = insertPhoneNumber(db, "(123) 456-7893")
	must(err)
	_, err = insertPhoneNumber(db, "123-456-7894")
	must(err)
	_, err = insertPhoneNumber(db, "123-456-7890")
	must(err)
	_, err = insertPhoneNumber(db, "1234567892")
	must(err)
	_, err = insertPhoneNumber(db, "(123) 456-7890")
	must(err)

	number, err := getPhoneNumber(db, id)
	must(err)
	fmt.Println("Number is :", number)

	phones, err := allPhoneNumber(db)
	must(err)
	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Number)
		if number != p.Number {
			existing, err := findPhone(db, number)
			must(err)
			if existing != nil {
				must(deletePhone(db, p.ID))
			} else {
				p.Number = number
				must(updatePhone(db, p))
			}
		} else {
			fmt.Println("no change required")
		}
	}

}

func must(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
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

func getPhoneNumber(db *sql.DB, id int) (string, error) {
	var number string
	row := db.QueryRow("SELECT phone_number FROM phone_numbers WHERE id = $1", id)
	err := row.Scan(&number)
	if err != nil {
		return "", err
	}
	return number, nil

}

func allPhoneNumber(db *sql.DB) ([]phone, error) {
	rows, err := db.Query("SELECT id, phone_number FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phoneNumbers []phone
	for rows.Next() {
		var p phone
		err := rows.Scan(&p.ID, &p.Number)
		if err != nil {
			return nil, err
		}
		phoneNumbers = append(phoneNumbers, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return phoneNumbers, nil
}

func findPhone(db *sql.DB, number string) (*phone, error) {
	var p phone
	row := db.QueryRow("SELECT * FROM phone_numbers WHERE phone_number = $1", number)
	err := row.Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err

	}
	return &p, nil
}

func deletePhone(db *sql.DB, id int) error {
	statement := `DELETE FROM phone_numbers WHERE id = $1`
	_, err := db.Exec(statement, id)
	return err
}

func updatePhone(db *sql.DB, p phone) error {
	statement := `UPDATE phone_numbers SET phone_number = $2 WHERE id = $1`
	_, err := db.Exec(statement, p.ID, p.Number)
	return err
}

// func normalize(phone string) string {
// 	var buf bytes.Buffer

// 	// when iterate, string will be []rune
// 	for _, ch := range phone {
// 		// condition to check if ch >= 0 && ch <= 9
// 		if ch >= '0' && ch <= '9' {§
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
