package model

import "time"

var ErrorCodes = map[int]string{
	400: "bad_request",
	401: "unauthorized",
	403: "forbidden",
	404: "not_found",
	409: "conflict",
	422: "unprocessable_entity",
	429: "too_many_requests",
	500: "internal_error",
}

// getErrorCode returns a human-readable description for an error code.
func getErrorCode(status int) string {
	if rc, ok := ErrorCodes[status]; ok {
		return rc
	}
	return "unknown_error_code"
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code      string       `json:"code"`
	Status    int          `json:"status"`
	Detail    *string      `json:"detail,omitempty"`
	Instance  *string      `json:"instance,omitempty"`
	Timestamp *time.Time   `json:"timestamp,omitempty"`
	Errors    []FieldError `json:"errors,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func GetErrorResponse(status int, detail string, instance string, errors []FieldError) ErrorResponse {
	return ErrorResponse{
		Error: Error{
			Code:      getErrorCode(status),
			Status:    status,
			Detail:    new(detail),
			Instance:  new(instance),
			Timestamp: new(time.Now()),
			Errors:    errors,
		},
	}
}
