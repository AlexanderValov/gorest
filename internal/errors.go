package internal

import "errors"

var ErrEmptyUsername = errors.New("the username is required")

var ErrUsernameExist = errors.New("the username is already exist")

var ErrInvalidSizeOfUsername = errors.New("the username must be between 2-40 chars")

var ErrNoUser = errors.New("no user")
