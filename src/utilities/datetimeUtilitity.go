package utilities

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func GetNextMonth(date string) (string, error) {
	var result string
	dateSplit := strings.Split(date, "-")
	if len(dateSplit) != 3 {
		return result, errors.New("Invalid format..")
	}
	month, _ := strconv.ParseUint(dateSplit[1], 10, 16)
	year, _ := strconv.ParseUint(dateSplit[0], 10, 16)

	if month+1 > 12 {
		month = 1
		year++
	} else {
		month++
	}

	result = fmt.Sprintf("%d-%02d-%s", year, month, dateSplit[2])
	return result, nil
}
