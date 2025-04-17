package store

import (
	"log"
	"strconv"
	"time"
)

type Person struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      int    `json:"phone"`
	LastAccess string `json:"updated_at"`
}

func initPersonEntry(name, last_name string, number int) *Person {
	if name == "" || last_name == "" {
		return nil
	}
	// Give LastAccess a value
	LastAccess := strconv.FormatInt(time.Now().Unix(), 10)
	return &Person{name, last_name, number, LastAccess}
}

type DB interface {
	CountRecords() int
	SearchStartWith(number int) []Person
	Search(number int) *Person
	Remove(phone int) error
	Insert(first_name string, last_name string, phone int) error
	List() ([]Person, error)
	// updateIndexes(listRecords []Person) error
	initDb() error
}

func GetDB() DB {
	db := &PostgresDb{}
	err := db.initDb()
	if err != nil {
		log.Println(err)
	}
	return db
}
