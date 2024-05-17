package middleware

import (
  "net/http"
  "time"
  "github.com/gin-gonic/gin"
  "cems-dis/model"
  rs "cems-dis/server/response"
  "cems-dis/utils"
)


func reject(c *gin.Context, resp rs.Response) {
  resp.Json(c)
  c.Abort()
}

func TokenAuthMiddleware(c *gin.Context) {
  mdl, err := model.New()
  if err != nil {
    panic(err.Error())
  }
  
  token := utils.ParseBearerToken(c.GetHeader("Authorization"))
  if len(token) == 0 {
    reject(c, rs.Error(http.StatusBadRequest, "Missing bearer token"))
    return
  }
  
  deviceToken, err := mdl.GetDeviceToken(token)
	if err != nil {
    reject(c, rs.Error(http.StatusInternalServerError, "DB error"))
    return
  } else if deviceToken == nil {
    reject(c, rs.Error(http.StatusBadRequest, "Akses token tidak valid"))
    return
  } else if deviceToken.RefreshExpiredAt.Before(time.Now()) {
		reject(c, rs.Error(http.StatusBadRequest, "Akses token expired"))
    return
	}

  c.Next()
}
