package main

import (
  "fmt"
  "time"

  "os"
  "os/signal"
  "syscall"
  log "github.com/sirupsen/logrus"
  "cems-dis/config"
  "cems-dis/internal/token_store"
  "cems-dis/internal/transmitter"
  "cems-dis/internal/workers"
  "cems-dis/model"
  "cems-dis/server"
)

func init_xxx() {
  c := make(chan os.Signal)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go func() {
    <-c
    fmt.Println("SIGTERM received")
    transmitter.Stop()
    workers.StopWorkersAndWait()
    time.Sleep(time.Duration(3 * time.Second))
    os.Exit(1)
  }()
}

func main() {
  config.Load("application")
  model, err := model.New()
  if err != nil {
    log.Warnf("Error connecting to database: %s", err.Error())
    os.Exit(1)
  }
  srv := server.New(model)
  workers.InitWorkers()
  token_store.Init()
  transmitter.Init(model)
  go transmitter.Start()
  srv.Start()
}
