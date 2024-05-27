package utils

import (
  "time"
)

const defaultDateTimeFormat = "2006-01-02 15:04:05"

func TimeToString(t time.Time) string {
  if t.IsZero() {
    return ""
  }
  return t.Format(defaultDateTimeFormat)
}
