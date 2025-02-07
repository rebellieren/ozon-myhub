package utils

import "fmt"

type GraphQLError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func NewGraphQLError(message, code string) *GraphQLError {
	return &GraphQLError{
		Message: message,
		Code:    code,
	}
}

func (e *GraphQLError) Error() string {
	return fmt.Sprintf("%s (code: %s)", e.Message, e.Code)
}
