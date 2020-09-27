package pkg

import "github.com/gobuffalo/validate/v3"

type MyError struct {
	*validate.Errors
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewMyError(err error, code int) *MyError {
	verr, ok := err.(*validate.Errors)
	if !ok {
		verr = validate.NewErrors()
	}

	msg := GetMessageFromErrorCodeMap(code)
	if msg == "" {
		msg = err.Error()
	}

	return &MyError{Errors: verr, Code: code, Message: msg}
}

func GetMessageFromErrorCodeMap(code int) string {
	if msg, ok := errorCodeMessageMap[code]; ok {
		return msg
	}
	return ""
}

// errorCodeMessageMap is a error code map manager
var errorCodeMessageMap = map[int]string{
	// 1000 - 2000 for user relevant error codes
	1000: "invalid password",
	1001: "incorrect login data",
	1002: "incorrect registration data",
	1003: "",
	1004: "incorrect data",

	// database error
	9999: "database error",
}
