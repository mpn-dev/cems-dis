package main

import (
  "os"
  log "github.com/sirupsen/logrus"
  "cems-dis/config"
  "cems-dis/model"
  "cems-dis/server"
)

func main() {
  config.Load("application")
  model, err := model.New()
  if err != nil {
    log.Warnf("Error connecting to database: %s", err.Error())
    os.Exit(1)
  }
  server.Start(model)
}
