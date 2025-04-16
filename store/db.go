package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/phone_book/lib"
)

type Person struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      int    `json:"phone"`
	LastAccess string `json:"updated_at"`
}

type DB struct {
	path string
	// map [phone]index
	indexes map[int]int
}

func (db *DB) initDb() {
	// inint store
	if _, err := os.Stat(db.path); err != nil && errors.Is(err, os.ErrNotExist) {
		errInitDb := lib.WriteSerializeJSONFile(db.path, make([]any, 0))
		if errInitDb != nil {
			fmt.Printf("init db error => %s", errInitDb)
			return
		}
	}
	// inint indexes
	if _, err := os.Stat("store/indexes.json"); err != nil && errors.Is(err, os.ErrNotExist) {
		errInitIndexes := lib.WriteSerializeJSONFile("store/indexes.json", make(map[int]int))
		if errInitIndexes != nil {
			fmt.Printf("init indexes error => %s", errInitIndexes)
			return
		}
	} else {
		errOpenIndx := lib.OpenDeSerializeJSONFile("store/indexes.json", &db.indexes)
		if errOpenIndx != nil {
			fmt.Printf("open indexes db error => %s", errOpenIndx)
			return
		}
	}
}

func (db *DB) updateIndexes(listRecords []Person) error {
	newIndexes := make(map[int]int)
	for i, r := range listRecords {
		newIndexes[r.Phone] = i
	}
	db.indexes = newIndexes
	return lib.WriteSerializeJSONFile("store/indexes.json", newIndexes)
}

func (db *DB) CountRecords() int {
	return len(db.indexes)
}

func initPersonEntry(name, last_name string, number int) *Person {
	if name == "" || last_name == "" {
		return nil
	}
	// Give LastAccess a value
	LastAccess := strconv.FormatInt(time.Now().Unix(), 10)
	return &Person{name, last_name, number, LastAccess}
}

func (db *DB) SearchStartWith(number int) []Person {
	searchNum := strconv.Itoa(number)
	listRecords, _ := db.List()
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

func (db *DB) Search(number int) *Person {
	_, ok := db.indexes[number]
	if !ok {
		return nil
	}

	listRecords, _ := db.List()
	return &listRecords[db.indexes[number]]

}

func (db *DB) Remove(phone int) error {
	listRecords, _ := db.List()
	i, ok := db.indexes[phone]
	if !ok {
		return fmt.Errorf("Record with number %v not found!", phone)
	}

	listRecords = append(listRecords[:i], listRecords[i+1:]...)
	err := db.updateIndexes(listRecords)
	if err != nil {
		return err
	}

	err = lib.WriteSerializeJSONFile(db.path, listRecords)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Insert(first_name string, last_name string, phone int) error {
	temp := initPersonEntry(first_name, last_name, phone)
	f, err := os.OpenFile(db.path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	listRecords, _ := db.List()

	_, ok := db.indexes[phone]
	if ok {
		return fmt.Errorf("Person with number: %d already exsist!", phone)
	}

	listRecords = append(listRecords, *temp)

	encodedListRecords, err := json.MarshalIndent(listRecords, "", "    ")
	if err != nil {
		return err
	}

	if _, err = f.Write(encodedListRecords); err != nil {
		return err
	}

	err = db.updateIndexes(listRecords)
	if err != nil {
		return err
	}

	return nil

}

func (db *DB) List() ([]Person, error) {
	f, err := os.Open(db.path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	// b, err := io.ReadAll(f)
	// if err != nil {
	// 	return nil, err
	// }

	outputSlice := []Person{}

	errDecode := lib.DeSerialize(&outputSlice, f)
	if errDecode != nil {
		return nil, errDecode
	}

	db.updateIndexes(outputSlice)

	return outputSlice, nil

}

func GetDB() DB {
	db := DB{"store/store.json", make(map[int]int)}
	db.initDb()
	return db
}
