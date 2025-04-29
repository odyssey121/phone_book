package store

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/phone_book/internal/config"
)

var (
	ErrPhoneExist = errors.New("phone exist")
)

type Person struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      int    `json:"phone"`
	LastAccess string `json:"updated_at"`
}

func InitPersonEntry(name, last_name string, number int) *Person {
	// Give LastAccess a value
	LastAccess := strconv.FormatInt(time.Now().Unix(), 10)
	return &Person{name, last_name, number, LastAccess}
}

const PostgresDriver = "postgres_driver"
const JsonDriver = "json_driver"

type DB interface {
	CountRecords() int
	SearchStartWith(number int) ([]Person, error)
	Search(number int) *Person
	Remove(phone int) error
	Insert(first_name string, last_name string, phone int) error
	List() ([]Person, error)
	init() error
}

func GetDB(cfg config.Storage) (DB, error) {
	var db DB
	switch cfg.Driver {
	case PostgresDriver:
		db = &PostgresDb{fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)}
	case JsonDriver:
		db = &JsonDb{Path: cfg.StoragePath, IndexesPath: cfg.IndexesPath}
	}

	err := db.init()
	if err != nil {
		return nil, err
	}
	return db, nil
}
