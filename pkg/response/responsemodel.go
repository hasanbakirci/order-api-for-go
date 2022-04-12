package response

import "github.com/labstack/echo/v4"

type ErrorDetails struct {
	StatusCode int         `json:"statusCode"`
	ErrorCode  int         `json:"errorCode"`
	Message    interface{} `json:"message"`
}

type SuccessDetails struct {
	Data       interface{} `json:"data"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
}

func CustomPanic(statusCode int, errorCode int, message string) ErrorDetails {
	return ErrorDetails{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func SuccessResponse(c echo.Context, statusCode int, data interface{}, message string) (err error) {
	successDetails := SuccessDetails{
		Data:       data,
		StatusCode: statusCode,
		Message:    message,
	}
	err = c.JSON(statusCode, successDetails)
	return
}

func ErrorResponse(c echo.Context, statusCode int, errorCode int, error interface{}) (err error) {
	errorDetails := ErrorDetails{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    error,
	}
	err = c.JSON(statusCode, errorDetails)
	return
}
