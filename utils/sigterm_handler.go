package utils

import (
  "os"
  "os/signal"
  "syscall"
)

func HandleSigTerm(cbk func()) {
  c := make(chan os.Signal)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go func() {
    <-c
    cbk()
    os.Exit(1)
  }()
}
