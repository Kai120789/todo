package utils

import (
	"net/http"
	"strconv"
	"strings"
)

func ParseIDFromURL(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}
