package main

import (
  "fmt"
  "os"

  log "github.com/sirupsen/logrus"
  "cems-dis/config"
  "cems-dis/internal/logging"
  "cems-dis/internal/token_store"
  "cems-dis/internal/transmitter"
  "cems-dis/internal/workers"
  "cems-dis/model"
  "cems-dis/server"
  "cems-dis/utils"
)

func init() {
  utils.HandleSigTerm(func() {
    fmt.Println("SIGTERM received")
    transmitter.Stop()
    workers.StopWorkersAndWait()
  })

  logging.MoveLogContent()
}

func main() {
  logging.Configure()
  config.Load("application")
  model, err := model.New()
  if err != nil {
    log.Warningf("Error connecting to database: %s", err.Error())
    os.Exit(1)
  }
  srv := server.New(model)
  workers.InitWorkers()
  token_store.Init()
  transmitter.Init(model)
  go transmitter.Start()
  srv.Start()
}
