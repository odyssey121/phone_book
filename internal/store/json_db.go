package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/phone_book/internal/lib"
)

type JsonDb struct {
	Path        string
	IndexesPath string
	indexes     map[int]int
}

func (db *JsonDb) init() error {
	// inint store
	db.indexes = make(map[int]int)

	if _, err := os.Stat(db.Path); err != nil && errors.Is(err, os.ErrNotExist) {
		errInitDb := lib.WriteSerializeJSONFile(db.Path, make([]any, 0))
		if errInitDb != nil {
			return fmt.Errorf("init db error => %s", errInitDb)
		}
	}
	// inint indexes
	if _, err := os.Stat(db.IndexesPath); err != nil && errors.Is(err, os.ErrNotExist) {
		errInitIndexes := lib.WriteSerializeJSONFile(db.IndexesPath, make(map[int]int))
		if errInitIndexes != nil {
			return fmt.Errorf("init indexes error => %s", errInitIndexes)
		}
	} else {
		errOpenIndx := lib.OpenDeSerializeJSONFile(db.IndexesPath, &db.indexes)
		if errOpenIndx != nil {
			return fmt.Errorf("open indexes db error => %s", errOpenIndx)
		}
	}
	return nil
}

func (db *JsonDb) updateIndexes(listRecords []Person) error {
	newIndexes := make(map[int]int)
	for i, r := range listRecords {
		newIndexes[r.Phone] = i
	}
	db.indexes = newIndexes
	return lib.WriteSerializeJSONFile(db.IndexesPath, newIndexes)
}

func (db *JsonDb) CountRecords() int {
	return len(db.indexes)
}

func (db *JsonDb) SearchStartWith(number int) []Person {
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

func (db *JsonDb) Search(number int) *Person {
	_, ok := db.indexes[number]
	if !ok {
		return nil
	}

	listRecords, _ := db.List()
	return &listRecords[db.indexes[number]]

}

func (db *JsonDb) Remove(phone int) error {
	listRecords, _ := db.List()
	i, ok := db.indexes[phone]
	if !ok {
		return fmt.Errorf("record with number %d not exist", phone)
	}

	listRecords = append(listRecords[:i], listRecords[i+1:]...)
	err := db.updateIndexes(listRecords)
	if err != nil {
		return err
	}

	err = lib.WriteSerializeJSONFile(db.Path, listRecords)

	if err != nil {
		return err
	}

	return nil
}

func (db *JsonDb) Insert(first_name string, last_name string, phone int) error {
	temp := initPersonEntry(first_name, last_name, phone)
	f, err := os.OpenFile(db.Path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	listRecords, _ := db.List()

	_, ok := db.indexes[phone]
	if ok {
		return fmt.Errorf("person with number: %d already exsist", phone)
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

func (db *JsonDb) List() ([]Person, error) {
	f, err := os.Open(db.Path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	outputSlice := []Person{}

	errDecode := lib.DeSerialize(&outputSlice, f)
	if errDecode != nil {
		return nil, errDecode
	}

	db.updateIndexes(outputSlice)

	return outputSlice, nil

}
