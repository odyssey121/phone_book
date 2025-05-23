package store

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresDb struct {
	ConnectionString string
}

func (db *PostgresDb) getConnect() *sql.DB {
	conn, err := sql.Open("postgres", db.ConnectionString)
	if err != nil {
		log.Println(err)
		return nil
	}

	return conn
}

func (db *PostgresDb) init() error {
	const op = "storage.postgres_db.init"
	conn := db.getConnect()
	defer conn.Close()

	stmt, err := conn.Prepare(
		`CREATE TABLE IF NOT EXISTS Persons (
			first_name varchar(255) NOT NULL,
			last_name varchar(255) NOT NULL,
			phone VARCHAR(22) NOT NULL UNIQUE,
			last_access integer DEFAULT NULL,
			PRIMARY KEY (phone)
		)`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (db *PostgresDb) CountRecords() int {
	var count int
	conn := db.getConnect()
	defer conn.Close()

	row := conn.QueryRow("SELECT count(*) FROM Persons")
	err := row.Scan(&count)
	if err != nil {
		log.Println("CountRecords() failed:", err)
	}

	return count
}

func (db *PostgresDb) SearchStartWith(number int) ([]Person, error) {
	const op = "storage.postgres_db.SearchStartWith"
	conn := db.getConnect()
	defer conn.Close()
	result := []Person{}
	stmtOut, err := conn.Prepare("SELECT first_name, last_name, phone, last_access FROM Persons WHERE phone LIKE $1")
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmtOut.Query(strconv.Itoa(number) + "%")
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var firstName, lastName, lastAccess string
		var phone int
		err = rows.Scan(&firstName, &lastName, &phone, &lastAccess)
		if err != nil {
			return result, fmt.Errorf("%s: %w", op, err)
		}
		p := Person{firstName, lastName, phone, lastAccess}
		result = append(result, p)
	}

	return result, nil

}

func (db *PostgresDb) Search(number int) *Person {
	conn := db.getConnect()
	defer conn.Close()
	var firstName, lastName, lastAccess string
	var phone int
	row := conn.QueryRow("SELECT first_name, last_name, phone, last_access FROM Persons WHERE phone = $1 LIMIT 1", number)
	row.Scan(&firstName, &lastName, &phone, &lastAccess)
	if phone == 0 {
		return nil
	}

	return &Person{firstName, lastName, phone, lastAccess}

}

func (db *PostgresDb) Remove(phone int) error {
	const op = "storage.postgres_db.Remove"
	conn := db.getConnect()
	defer conn.Close()
	stmt, err := conn.Prepare("DELETE FROM Persons WHERE phone = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(phone)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (db *PostgresDb) Insert(first_name string, last_name string, phone int) error {
	const op = "storage.postgres_db.Insert"
	conn := db.getConnect()
	defer conn.Close()

	stmt, err := conn.Prepare(`INSERT INTO Persons (first_name, last_name, phone, last_access) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(first_name, last_name, phone, time.Now().Unix())

	if err != nil {
		if psgErr, ok := err.(*pq.Error); ok && psgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, ErrPhoneExist)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (db *PostgresDb) List() ([]Person, error) {
	conn := db.getConnect()
	defer conn.Close()
	result := []Person{}
	rows, _ := conn.Query("SELECT * FROM Persons")
	for rows.Next() {
		var firstName, lastName, lastAccess string
		var phone int
		rows.Scan(&firstName, &lastName, &phone, &lastAccess)
		p := Person{firstName, lastName, phone, lastAccess}
		result = append(result, p)

	}

	return result, nil

}
