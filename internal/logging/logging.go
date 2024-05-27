package logging

import (
  "fmt"
  "io/ioutil"
  "os"
  "strings"
  "time"
  log "github.com/sirupsen/logrus"
  "cems-dis/utils"
)

type LogFormat struct {}

func (f *LogFormat) Format(entry *log.Entry) ([]byte, error) {
  msg := fmt.Sprintf(
    "%s [%s] %s\n", 
    utils.TimeToString(entry.Time), 
    strings.ToUpper(entry.Level.String()), 
    entry.Message)
  return []byte(msg), nil
}


func Configure() {
  if true {
    log.SetFormatter(&LogFormat{})
    file, err := os.OpenFile("logs/production.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err == nil {
    log.SetOutput(file)
    }
  }
}

func MoveLogContent() {
  go func() {
    for {
      now := time.Now()
      if (now.Unix() % 604800) == 0 {
        src := "logs/production.log"
        dst := fmt.Sprintf("%s.%s", src, now.Format("20060102150405"))
        _ = os.Remove(dst)
        content, err := ioutil.ReadFile(src)
        if err != nil {
          log.Warningf("Moving log content failed: Error reading log file")
          continue
        }
        if len(content) == 0 {
          continue
        }
        err = ioutil.WriteFile(dst, content, 0666)
        _ = ioutil.WriteFile(src, []byte{}, 0666)
        log.Infof("Log content moved to %s\n", dst)
        time.Sleep(time.Duration(5 * time.Second))
      }
    }
  }()
}
