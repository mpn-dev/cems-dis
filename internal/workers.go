package internal

import (
  "fmt"
  "math/rand"
  "time"
  "github.com/alitto/pond"
)

var pool pond.Pond

func StartWorker(fn func()) {
	pool.Submit(fn)
}

func StopWorkersAndWait() {
	pool.StopAndWait()
}

func RunningWorkers() int {
	return pool.RunningWorkers()
}

func InitWorkers() {
  pool = pond.New(100, 1000)
}
