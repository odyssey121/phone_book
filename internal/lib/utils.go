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
		return 0, fmt.Errorf("phone number \"%s\" is incorrect", num)
	}
	n, _ := strconv.Atoi(NumFormated)
	return n, nil
}

// Serialize serializes a slice with JSON records
func Serialize(slice interface{}, w io.Writer) error {
	const op = "utils.Serialize"
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	err := e.Encode(slice)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// PrettyPrintJSONstream pretty prints the contents of the phone book
func PrettyPrintJSONstream(data interface{}) (string, error) {
	const op = "utils.PrettyPrintJSONstream"
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return buffer.String(), nil
}

// DeSerialize decodes a serialized slice with JSON records
func DeSerialize(slice interface{}, r io.Reader) error {
	const op = "utils.DeSerialize"
	e := json.NewDecoder(r)
	err := e.Decode(slice)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func WriteSerializeJSONFile(path string, slice interface{}) error {
	const op = "utils.WriteSerializeJSONFile"
	fd, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer fd.Close()

	err = Serialize(&slice, fd)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func OpenDeSerializeJSONFile(path string, slice interface{}) error {
	const op = "utils.OpenDeSerializeJSONFile"
	fd, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer fd.Close()

	return DeSerialize(slice, fd)

}
