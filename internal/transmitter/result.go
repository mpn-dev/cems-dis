package transmitter

import (
	"fmt"
	"strings"
	"cems-dis/model"
)

type Result struct {
	task			*model.Transmission
	success		bool
	code			int
	note			string
}

func (r Result) Task() *model.Transmission {
	return r.task
}

func (r Result) IsSuccess() bool {
	return r.success
}

func (r Result) Code() int {
	return r.code
}

func (r Result) Note() string {
	return r.note
}

func (r Result) Info() string {
	info := []string{}
	if r.success {
		info = append(info, "status: Success")
	} else {
		info = append(info, "status: Error")
	}
	if r.task != nil {
		info = append(info, fmt.Sprintf("transmission.Id: %d, protocol: %s, url: %s", r.task.Id, r.task.Protocol, r.task.BaseURL))
	}
	if len(r.note) > 0 {
		info = append(info, fmt.Sprintf("note: %s", r.note))
	}
	return strings.Join(info, ", ")
}

func Success(task *model.Transmission, code int, note string) Result {
	return Result{
		task:			task, 
		success: 	true, 
		code:			code, 
		note:			note, 
	}
}

func Error(task *model.Transmission, code int, note string) Result {
	return Result{
		task:			task, 
		success: 	false, 
		code:			code, 
		note:			note, 
	}
}
