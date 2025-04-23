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
	const op = "storage.json_db.init"
	// init store
	db.indexes = make(map[int]int)

	if _, err := os.Stat(db.Path); err != nil && errors.Is(err, os.ErrNotExist) {
		err := lib.WriteSerializeJSONFile(db.Path, make([]any, 0))
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	// inint indexes
	if _, err := os.Stat(db.IndexesPath); err != nil && errors.Is(err, os.ErrNotExist) {
		err := lib.WriteSerializeJSONFile(db.IndexesPath, make(map[int]int))
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	} else {
		err := lib.OpenDeSerializeJSONFile(db.IndexesPath, &db.indexes)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
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

func (db *JsonDb) SearchStartWith(number int) ([]Person, error) {
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
	return findedRecords, nil

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
	const op = "storage.json_db.Remove"
	listRecords, err := db.List()
	if err != nil {
		return fmt.Errorf("%s: %w", "storage.json_db.Remove.List", err)
	}
	i, ok := db.indexes[phone]
	if !ok {
		return nil
	}

	listRecords = append(listRecords[:i], listRecords[i+1:]...)
	err = db.updateIndexes(listRecords)
	if err != nil {
		return fmt.Errorf("%s: %w", "storage.json_db.Remove.updateIndexes", err)
	}

	err = lib.WriteSerializeJSONFile(db.Path, listRecords)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (db *JsonDb) Insert(first_name string, last_name string, phone int) error {
	const op = "storage.json_db.Insert"
	temp := initPersonEntry(first_name, last_name, phone)
	f, err := os.OpenFile(db.Path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	listRecords, _ := db.List()
	_, ok := db.indexes[phone]
	if ok {
		return fmt.Errorf("%s: %w", op, ErrPhoneExist)
	}

	listRecords = append(listRecords, *temp)

	encodedListRecords, err := json.MarshalIndent(listRecords, "", "    ")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = f.Write(encodedListRecords); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = db.updateIndexes(listRecords)
	if err != nil {
		return fmt.Errorf("%s: %w", "storage.json_db.Insert.updateIndexes", err)
	}

	return nil

}

func (db *JsonDb) List() ([]Person, error) {
	const op = "storage.json_db.List"
	f, err := os.Open(db.Path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	outputSlice := []Person{}

	err = lib.DeSerialize(&outputSlice, f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.updateIndexes(outputSlice)
	if err != nil {
		return outputSlice, fmt.Errorf("%s: %w", "storage.json_db.List.updateIndexes", err)
	}

	return outputSlice, nil

}
