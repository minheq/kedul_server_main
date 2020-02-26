package app

import (
	"fmt"
	"strings"
)

func makeIDsArgs(ids []string) (string, []interface{}) {
	args := make([]interface{}, len(ids))

	for i, id := range ids {
		args[i] = id
	}

	params := make([]string, 0, len(ids))

	for i := range ids {
		params = append(params, fmt.Sprintf("$%d", i+1))
	}

	return strings.Join(params, ", "), args
}
