package handler

import (
	"github.com/MayukhSobo/scaffold/internal/repository/users"
	"github.com/MayukhSobo/scaffold/pkg/http"
)

// ToUserResponse converts a database User model to a response-safe format using redaction
func ToUserResponse(user *users.User) *users.User {
	// Create a copy of the user to avoid modifying the original
	userCopy := *user

	// Redact sensitive fields marked with redact:"true" tag
	http.Redact(&userCopy)

	return &userCopy
}

// ToUserResponses converts a slice of database User models to response-safe format
func ToUserResponses(userSlice []users.User) []*users.User {
	responses := make([]*users.User, len(userSlice))
	for i, user := range userSlice {
		responses[i] = ToUserResponse(&user)
	}
	return responses
}
