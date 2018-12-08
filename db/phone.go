package db

import (
	"database/sql"
	"fmt"
)

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
