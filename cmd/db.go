package cmd

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strconv"
)

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     int    `json:"phone"`
}

type DB struct {
	path string
}

func (db *DB) searchStartWith(number int) []Person {
	searchNum := strconv.Itoa(number)
	listRecords, _ := db.list()
	findedRecords := []Person{}
	re := regexp.MustCompile(`^` + regexp.QuoteMeta(searchNum) + `.*`)
	for _, record := range listRecords {
		numStr := strconv.Itoa(record.Phone)
		if re.Match([]byte(numStr)) {
			findedRecords = append(findedRecords, record)
		}
	}
	return findedRecords

}

func (db *DB) search(number int) *Person {
	listRecords, _ := db.list()
	for _, record := range listRecords {
		if record.Phone == number {
			return &record
		}
	}
	return nil

}

func (db *DB) remove(phone int) error {
	listRecords, _ := db.list()
	for i, record := range listRecords {
		if record.Phone == phone {
			listRecords = append(listRecords[:i], listRecords[i+1:]...)
			break
		}
	}
	encodedListRecords, err := json.Marshal(listRecords)
	if err != nil {
		return err
	}

	if err = os.WriteFile(db.path, encodedListRecords, 0600); err != nil {
		return err
	}
	return nil
}

func (db *DB) insert(first_name string, last_name string, phone int) error {
	temp := Person{first_name, last_name, phone}
	f, err := os.OpenFile(db.path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	listRecords, _ := db.list()
	listRecords = append(listRecords, temp)

	encodedListRecords, err := json.Marshal(listRecords)
	if err != nil {
		return err
	}

	if _, err = f.Write(encodedListRecords); err != nil {
		return err
	}

	return nil

}

func (db *DB) list() ([]Person, error) {
	f, err := os.Open(db.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	v := []Person{}
	json.Unmarshal(b, &v)
	return v, nil

}

func getDB() DB {
	return DB{"db/phones.json"}
}
