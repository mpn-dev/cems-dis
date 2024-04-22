package config

import (
  "fmt"
  "os"
  "testing"

  "github.com/stretchr/testify/assert"
)


func TestConfig(t *testing.T) {
  vars := map[string]string{
    "SERVER_PORT":                  "2000", 
    "DB_HOST":                      "localhost", 
    "DB_PORT":                      "5432", 
    "DB_USER":                      "test_user", 
    "DB_PASS":                      "test_pass", 
    "DB_NAME":                      "test_name", 
    "USER_LOGIN_TOKEN_AGE_SHORT":   "1800", 
    "USER_LOGIN_TOKEN_AGE_LONG":    "86400", 
    "DEVICE_LOGIN_TOKEN_AGE":       "310", 
    "DEVICE_REFRESH_TOKEN_AGE":     "320", 
  }

  for k, v := range vars {
    os.Setenv(fmt.Sprintf("CEMS_%s", k), v)
    defer os.Unsetenv(k)
  }

  Load("application")

  assert.Equal(t, 2000, ServerPort())
  assert.Equal(t, "postgres://test_user:test_pass@localhost:5432/test_name", DbConfig().String())
  assert.Equal(t, 1800, UserLoginTokenAgeShort())
  assert.Equal(t, 86400, UserLoginTokenAgeLong())
  assert.Equal(t, 310, DeviceLoginTokenAge())
  assert.Equal(t, 320, DeviceRefreshTokenAge())
}
