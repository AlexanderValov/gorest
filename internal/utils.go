package internal

import "errors"

func GetErrResponse(err error) map[string]string {
	response := make(map[string]string, 1)
	switch {
	case errors.Is(err, ErrUsernameExist):
		response["ERR"] = ErrUsernameExist.Error()
	case errors.Is(err, ErrEmptyUsername):
		response["ERR"] = ErrEmptyUsername.Error()
	case errors.Is(err, ErrInvalidSizeOfUsername):
		response["ERR"] = ErrInvalidSizeOfUsername.Error()
	case errors.Is(err, ErrNoUser):
		response["ERR"] = ErrNoUser.Error()
	default:
		response["ERR"] = "unknown error, check logs"
	}
	return response
}
