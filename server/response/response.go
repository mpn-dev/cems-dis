package response

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

type Response struct {
  Code      int
  Error     interface{}
  Data      interface{}
  Formats   []Formatter
}

func (r Response) IsSuccess() bool {
  return (r.Code >= 200) && (r.Code <= 299)
}

func (r Response) IsError() bool {
  return r.Code != http.StatusOK
}

func (r Response) FormatWith(f ...Formatter) Response {
  r.Formats = f
  return r
}

func (r Response) Json(c *gin.Context) {
  ff := r.Formats[0:len(r.Formats)]
  if len(ff) == 0 {
    ff = append(ff, NewDefaultFormatter())
  }
  data := gin.H{}
  for _, f := range ff {
    data = f.H(r, data)
  }
  c.JSON(r.Code, data)
}

func New(code int, error interface{}, data interface{}) Response {
  return Response{
    Code:     code, 
    Error:    error, 
    Data:     data, 
  }
}

func Success(data interface{}) Response {
  return New(http.StatusOK, nil, data)
}

func Error(code int, message string) Response {
  return New(code, message, nil)
}
