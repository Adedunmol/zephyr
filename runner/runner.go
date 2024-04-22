package runner

import (
	"time"
)

type Task interface {
	Destroy()
}

type Runner struct {
	tasks   chan Task
	timeout <-chan time.Time
}

func New(d time.Duration) *Runner {

	return &Runner{
		tasks:   make(chan Task, 10),
		timeout: time.After(d),
	}
}
