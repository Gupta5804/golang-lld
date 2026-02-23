package main

import "sync"

type Job struct {
	ID       int
	FilePath string
}
type Result struct {
	JobID  int
	Status string
}

type Dispatcher struct {
	workerCount int
	jobs        chan Job
	wg          sync.WaitGroup
}

func NewDispatcher(workerCount, bufferSize int) *Dispatcher {
	return &Dispatcher{
		jobs:        make(chan Job, bufferSize),
		workerCount: workerCount,
	}
}

func (d *Dispatcher) Start(results chan<- Result) {
	for i := 0; i < d.workerCount; i++ {
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			for job := range d.jobs {
				results <- Result{JobID: job.ID, Status: "done"}
			}
		}()
	}
}
func (d *Dispatcher) Publish(job Job) {
	d.jobs <- job
}
func (d *Dispatcher) Stop() {
	close(d.jobs)
	d.wg.Wait()
}
