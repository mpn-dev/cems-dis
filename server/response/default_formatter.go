package response

import "github.com/gin-gonic/gin"

type DefaultFormatter struct {}

func (f DefaultFormatter) Json(c *gin.Context, r Response) {
	c.JSON(r.Code, gin.H{"error": r.Error, "data": r.Data})
}

func NewDefaultFormatter() DefaultFormatter {
	return DefaultFormatter{}
}
