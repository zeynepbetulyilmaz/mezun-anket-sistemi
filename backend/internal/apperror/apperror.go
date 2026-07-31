// Package apperror: tüm handler'ların döndüğü tek tip hata yapısı.
// Middleware bunu HTTP status + standart JSON zarfına çevirir.
package apperror

import "net/http"

type Detail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	HTTPStatus int      `json:"-"`
	Code       string   `json:"code"`
	Message    string   `json:"message"`
	Details    []Detail `json:"details,omitempty"`
}

func (e *AppError) Error() string { return e.Message }

func Validation(message string, details ...Detail) *AppError {
	return &AppError{HTTPStatus: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: message, Details: details}
}

func Unauthorized(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: message}
}

func Forbidden(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusForbidden, Code: "FORBIDDEN", Message: message}
}

func NotFound(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusNotFound, Code: "NOT_FOUND", Message: message}
}

func Conflict(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusConflict, Code: "CONFLICT", Message: message}
}

func Internal(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: message}
}
