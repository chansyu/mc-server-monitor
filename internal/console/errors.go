package console

import "errors"

var ErrNoPlayer = errors.New("No player was found")
var ErrMalformedOutput = errors.New("Recieved malformed output")
