package utils

import (
	"strings"
)

func ParseBearerToken(auth string) string {
  bearer := strings.Trim(auth, " ")
  token := ""
  if len(bearer) > 7 {
    token = strings.Trim(bearer[7:len(bearer)], " ")
  }
  if strings.ToLower(strings.Trim(bearer[0:7], " ")) != "bearer" {
    return ""
  }
	return token
}
