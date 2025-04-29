package response

import "github.com/phone_book/internal/store"

type Response struct {
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type ResponseDataPersons struct {
	Data   []store.Person
	Error  string
	Status string
}

type ResponseDataPerson struct {
	Data   store.Person
	Error  string
	Status string
}

type ResponseData struct {
	Response
	Data interface{} `json:"data,omitempty"`
}

const (
	StatusOK               = "OK"
	StatusError            = "Error"
	InternalServerErrorMsg = "internal server error"
)

func OK(msg string) Response {
	return Response{Status: StatusOK, Message: msg}
}

func WithData(data interface{}) ResponseData {
	return ResponseData{
		Response{Status: StatusOK},
		data,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
