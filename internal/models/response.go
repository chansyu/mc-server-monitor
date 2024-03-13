package models

import "time"

const MsgError = "An error occurred."

type Response struct {
	Success   bool
	Message   string
	Timestamp time.Time
	Command   string
	Args      []string
}

func NewResponse(cmd string, arg []string) Response {
	return Response{
		Success:   false,
		Message:   "",
		Timestamp: time.Now(),
		Command:   cmd,
		Args:      arg,
	}
}

func (r *Response) Error() {
	r.Message = MsgError
}

func (r *Response) Succeed(msg string) {
	r.Message = msg
	r.Success = true
}
