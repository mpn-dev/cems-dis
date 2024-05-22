package response

import (
	"github.com/gin-gonic/gin"
)

type DefaultFormatter struct {}

func (f DefaultFormatter) H(r Response, d gin.H) gin.H {
	status := "success"
	if !r.IsSuccess() {
		status = "error"
	}

	data := addMetaField(d)
	data["meta"].(gin.H)["status"] 	= status
	data["meta"].(gin.H)["message"]	= r.Error

	if r.Data != nil {
		data["data"] = r.Data
	}

	return data
}

func NewDefaultFormatter() DefaultFormatter {
	return DefaultFormatter{}
}
