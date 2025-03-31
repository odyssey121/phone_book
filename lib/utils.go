package lib

import (
	"fmt"
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
