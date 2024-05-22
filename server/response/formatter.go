package response

import (
  "github.com/gin-gonic/gin"
)

type Formatter interface {
  H(Response, gin.H) gin.H
}


func addMetaField(d gin.H) gin.H {
  if _, ok := d["meta"]; !ok {
    d["meta"] = gin.H{}
  }
  return d
}
