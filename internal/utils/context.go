// utils/context.go
package utils

import (
	"context"
	"net/http"
)

type contextKey string

const userIDKey = contextKey("userID")

func GetUserIDFromContext(r *http.Request) int {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		return 0
	}
	return userID
}

func AddUserIDToContext(r *http.Request, userID int) *http.Request {
	ctx := context.WithValue(r.Context(), userIDKey, userID)
	return r.WithContext(ctx)
}
