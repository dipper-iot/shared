package worker_file

import "fmt"

type ErrorPayload struct {
	Code string
	Err  error
}

func NewErrorPayload(code string, err error) *ErrorPayload {
	return &ErrorPayload{Code: code, Err: err}
}

func (e ErrorPayload) Error() string {
	return fmt.Sprintf("Code: %s \n error: %s", e.Code, e.Err.Error())
}

func ConvertError(err error) (*ErrorPayload, error) {
	switch err.(type) {
	case ErrorPayload:
		data := err.(ErrorPayload)
		return &data, err
	case *ErrorPayload:
		data := err.(*ErrorPayload)
		return data, err
	default:
		return nil, err
	}
}
