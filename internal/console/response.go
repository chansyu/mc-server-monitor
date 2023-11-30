package console

import "time"

type Response struct {
	Success   bool
	Message   string
	Timestamp time.Time
	Command   string
	Args      []string
}

func newResponse(cmd string, arg []string) *Response {
	return &Response{
		Success:   false,
		Message:   "An error occurred",
		Timestamp: time.Now(),
		Command:   cmd,
		Args:      arg,
	}
}

func (r *Response) consoleDisconnect() {
	r.Message = "Request Disconnect"
}

func (r *Response) consoleSuccess(msg string) {
	r.Message = msg
	r.Success = true
}
