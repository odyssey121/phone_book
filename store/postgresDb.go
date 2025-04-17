package store

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDb struct {
	connectionString string
}

var (
	Hostname = "localhost"
	Port     = 5432
	Username = "username"
	Password = "pass"
	Database = "phone_book"
)

func (db *PostgresDb) getConnect() *sql.DB {
	conn, err := sql.Open("postgres", db.connectionString)
	if err != nil {
		log.Println(err)
		return nil
	}

	return conn
}

func (db *PostgresDb) initDb() error {
	db.connectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, Username, Password, Database)

	conn := db.getConnect()
	defer conn.Close()

	stmt, err := conn.Prepare(
		`CREATE TABLE IF NOT EXISTS Persons (
			first_name varchar(255) NOT NULL,
			last_name varchar(255) NOT NULL,
			phone VARCHAR(22) NOT NULL,
			last_access integer DEFAULT NULL,
			PRIMARY KEY (phone)
		)`)
	if err != nil {
		return fmt.Errorf("init db Prepare error => %s", err)
	}

	res, err := stmt.Exec()
	if err != nil {
		return fmt.Errorf("init db Exec error => %s", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("init db RowsAffected error => %s", err)
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

func (db *PostgresDb) SearchStartWith(number int) []Person {
	conn := db.getConnect()
	defer conn.Close()
	result := []Person{}
	stmtOut, err := conn.Prepare("SELECT first_name, last_name, phone, last_access FROM Persons WHERE phone LIKE $1")
	if err != nil {
		log.Printf("SearchStartWith Prepare db error => %s", err)
	}

	rows, err := stmtOut.Query(strconv.Itoa(number) + "%")
	if err != nil {
		log.Printf("SearchStartWith Query db error => %s", err)
	}

	for rows.Next() {
		var firstName, lastName, lastAccess string
		var phone int
		err = rows.Scan(&firstName, &lastName, &phone, &lastAccess)
		if err != nil {
			log.Printf("SearchStartWith Scan db error => %s", err)
		}
		p := Person{firstName, lastName, phone, lastAccess}
		result = append(result, p)
	}

	return result

}

func (db *PostgresDb) Search(number int) *Person {
	conn := db.getConnect()
	defer conn.Close()
	var firstName, lastName, lastAccess string
	var phone int
	row := conn.QueryRow("SELECT first_name, last_name, phone, last_access FROM Persons WHERE phone = $1", number)
	row.Scan(&firstName, &lastName, &phone, &lastAccess)
	if phone == 0 {
		return nil
	}

	return &Person{firstName, lastName, phone, lastAccess}

}

func (db *PostgresDb) Remove(phone int) error {
	conn := db.getConnect()
	defer conn.Close()
	stmt, err := conn.Prepare("DELETE FROM Persons WHERE phone = $1")
	if err != nil {
		return fmt.Errorf("Remove Prepare db error => %s", err)
	}

	res, err := stmt.Exec(phone)

	if err != nil {
		return fmt.Errorf("Remove Exec db error => %s", err)
	}

	affectedN, _ := res.RowsAffected()

	log.Printf("Remove db affected: %d with data[phone: %d]", affectedN, phone)

	return nil
}

func (db *PostgresDb) Insert(first_name string, last_name string, phone int) error {
	conn := db.getConnect()
	defer conn.Close()

	stmt, err := conn.Prepare(`INSERT INTO Persons (first_name, last_name, phone, last_access) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return fmt.Errorf("Insert Prepare db error => %s", err)
	}

	res, err := stmt.Exec(first_name, last_name, phone, time.Now().Unix())

	if err != nil {
		return fmt.Errorf("Insert Exec db error => %s", err)
	}

	affectedN, _ := res.RowsAffected()

	log.Printf("Insert db affected: %d with data[last_name: %s,] [first_name: %s], [phone: %d]", affectedN, last_name, first_name, phone)

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
