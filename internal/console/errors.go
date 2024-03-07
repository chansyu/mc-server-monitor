package console

import "errors"

var ErrNoPlayer = errors.New("console: no player was found")
var ErrMalformedOutput = errors.New("console: recieved malformed output")
var ErrTimeout = errors.New("console: connection timeout")
var ErrInternal = errors.New("console: internal error")
