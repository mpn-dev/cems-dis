package response

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

type Formatter interface {
  Json(c *gin.Context, r Response)
}

type Response struct {
  Code      int
  Error     interface{}
  Data      interface{}
  fmt       Formatter
}

func (r Response) IsSuccess() bool {
  return r.Code == http.StatusOK
}

func (r Response) IsError() bool {
  return r.Code != http.StatusOK
}

func (r Response) UseFormatter(f Formatter) Response {
  r.fmt = f
  return r
}

func (r Response) Json(c *gin.Context) {
  r.fmt.Json(c, r)
}

func New(code int, error interface{}, data interface{}) Response {
  return Response{
    Code:   code, 
    Error:  error, 
    Data:   data, 
    fmt:    NewDefaultFormatter(), 
  }
}

func Success(data interface{}) Response {
  return New(http.StatusOK, nil, data)
}

func Error(code int, message string) Response {
  return New(code, message, nil)
}
