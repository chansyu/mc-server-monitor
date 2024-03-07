package models

import "time"

const MsgError = "An error occurred"
const MsgDisconnect = "Request Disconnect"

type Response struct {
	Success   bool
	Message   string
	Timestamp time.Time
	Command   string
	Args      []string
}

func NewResponse(cmd string, arg []string) *Response {
	return &Response{
		Success:   false,
		Message:   "",
		Timestamp: time.Now(),
		Command:   cmd,
		Args:      arg,
	}
}

func (r *Response) ConsoleError() {
	r.Message = MsgError
}

func (r *Response) ConsoleDisconnect() {
	r.Message = MsgDisconnect
}

func (r *Response) ConsoleSuccess(msg string) {
	r.Message = msg
	r.Success = true
}
