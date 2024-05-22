package response

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type DasFormatter struct {}

func (f DasFormatter) H(r Response, d gin.H) gin.H {
	status := "success"
	if !r.IsSuccess() {
		status = "error"
	}

	errlist := []interface{}{}
	if r.Error != nil {
		errlist	= []interface{}{r.Error}
	}

	data := addMetaField(d)
	data["meta"].(gin.H)["status"] 	= status
	data["meta"].(gin.H)["message"] =	http.StatusText(r.Code)
	data["meta"].(gin.H)["errors"]	= errlist

	if r.Data != nil {
		data["data"] = r.Data
	}

	return data
}

func NewDasFormatter() DasFormatter {
	return DasFormatter{}
}

func (r Response) UseDasFormatter() Response {
	return r.FormatWith(NewDasFormatter())
}

func (r Response) UseDasFormatterWithPaging(paging *Paging) Response {
	r.Formats = []Formatter{
		NewDasFormatter(), 
		paging, 
	}
	return r
}
