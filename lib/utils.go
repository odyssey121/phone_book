package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func matchNumber(num string) bool {
	re := regexp.MustCompile(`^\d+$`)
	return re.Match([]byte(num))
}

func FormatNumber(num string) (int, error) {
	NumFormated := strings.ReplaceAll(num, "-", "")
	if !matchNumber(NumFormated) {
		return 0, fmt.Errorf("Phone Number \"%s\" is Incorrect!\n", num)
	}
	n, _ := strconv.Atoi(NumFormated)
	return n, nil
}

// Serialize serializes a slice with JSON records
func Serialize(slice interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	return e.Encode(slice)
}

// PrettyPrintJSONstream pretty prints the contents of the phone book
func PrettyPrintJSONstream(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// DeSerialize decodes a serialized slice with JSON records
func DeSerialize(slice interface{}, r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(slice)
}

func WriteSerializeJSONFile(path string, slice interface{}) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	err = Serialize(&slice, fd)
	if err != nil {
		return err
	}

	return nil

}

func OpenDeSerializeJSONFile(path string, slice interface{}) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	return DeSerialize(slice, fd)

}
