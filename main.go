package main

import (
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

func main() {
	// reset db first
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	must(db.Reset("postgres", psqlInfo, dbName))

	// connect to new database that has been crated
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbName)
	must(db.Migrate("postgres", psqlInfo))

	db, err := db.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	err = db.Seed()
	must(err)

	phones, err := db.AllPhones()
	must(err)
	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Number)
		if number != p.Number {
			fmt.Println("Updating or removeing : ", p.Number)
			existing, err := db.FindPhone(number)
			must(err)
			if existing != nil {
				must(db.DeletePhone(p.ID))
			} else {
				p.Number = number
				must(db.UpdatePhone(&p))
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
