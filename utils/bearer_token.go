package utils

import (
	"strings"
)

func ParseBearerToken(auth string) string {
  bearer := strings.Trim(auth, " ")
  if len(bearer) <= 7 {
    return ""
  }
  if strings.ToLower(strings.Trim(bearer[0:7], " ")) != "bearer" {
    return ""
  }
  return strings.Trim(bearer[7:len(bearer)], " ")
}
