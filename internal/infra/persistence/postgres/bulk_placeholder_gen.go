package postgres

import (
	"errors"
	"fmt"
	"strings"
)

func GenBulkPlaceholders(query string, args any, rowLength int) (string, error) {
	placeholders := []string{}
	placeholdersCounter := 1

	argsArray, ok := args.([]any)
	if !ok {
		return "", errors.New("array is not interface{} slice")
	}

	if len(argsArray)%rowLength != 0 {
		return "", errors.New("remainder of row length from array length is not zero")
	}

	for i := 0; i < len(argsArray)/rowLength; i++ {
		row := make([]string, rowLength)

		for i := 0; i < rowLength; i++ {
			row[i] = fmt.Sprintf("$%d", placeholdersCounter)
			placeholdersCounter += 1
		}

		placeholders = append(placeholders, strings.Replace("(?)", "?", strings.Join(row, ", "), 1))
	}

	return strings.Replace(query, "?", strings.Join(placeholders, ", "), 1), nil
}
