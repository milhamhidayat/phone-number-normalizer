package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	DBPostgre *sql.DB
}

type Phone struct {
	ID     int
	Number string
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{DBPostgre: db}, nil
}

func (db *DB) Close() error {
	return db.DBPostgre.Close()
}

func (db *DB) Seed() error {

	data := []string{"123456789",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123) 456-7890"}

	for _, number := range data {
		if _, err := insertPhoneNumber(db.DBPostgre, number); err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) AllPhones() ([]Phone, error) {
	rows, err := db.DBPostgre.Query("SELECT id, phone_number FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phoneNumbers []Phone
	for rows.Next() {
		var p Phone
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

func insertPhoneNumber(db *sql.DB, phone string) (int, error) {
	var id int
	statement := `INSERT INTO phone_numbers (phone_number) VALUES($1) RETURNING id`
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = createPhoneNumbersTable(db)
	if err != nil {
		return err
	}
	return db.Close()
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
	return createDB(db, name)
}

// func findPhone(db *sql.DB, number string) (*phone, error) {
// 	var p phone
// 	row := db.QueryRow("SELECT * FROM phone_numbers WHERE phone_number = $1", number)
// 	err := row.Scan(&p.ID, &p.Number)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		}
// 		return nil, err

// 	}
// 	return &p, nil
// }

func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	row := db.DBPostgre.QueryRow("SELECT * FROM phone_numbers WHERE phone_number = $1", number)
	err := row.Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err

	}
	return &p, nil
}

func (db *DB) DeletePhone(id int) error {
	statement := `DELETE FROM phone_numbers WHERE id = $1`
	_, err := db.DBPostgre.Exec(statement, id)
	return err
}

func (db *DB) UpdatePhone(p *Phone) error {
	statement := `UPDATE phone_numbers SET phone_number = $2 WHERE id = $1`
	_, err := db.DBPostgre.Exec(statement, p.ID, p.Number)
	return err
}

// func getPhoneNumber(db *sql.DB, id int) (string, error) {
// 	var number string
// 	row := db.DBPostgre.QueryRow("SELECT phone_number FROM phone_numbers WHERE id = $1", id)
// 	err := row.Scan(&number)
// 	if err != nil {
// 		return "", err
// 	}
// 	return number, nil

// }
