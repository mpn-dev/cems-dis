package workers

import (
  "github.com/alitto/pond"
)

var pool *pond.WorkerPool

func StartJob(fn func()) {
	pool.Submit(fn)
}

func StopWorkersAndWait() {
  if pool != nil {
  	pool.StopAndWait()
  }
}

func RunningWorkers() int {
	return pool.RunningWorkers()
}

func InitWorkers() {
  if pool == nil {
    pool = pond.New(100, 1000)
  }
}
