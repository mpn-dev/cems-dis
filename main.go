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
    log.Println("Error connecting to databaes")
    os.Exit(1)
  }
  server.Start(model)
}
