package workerpool

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	Req        *http.Request
	ResultChan chan<- *http.Response
	ErrorChan  chan<- error
}

type WorkerPool struct {
	jobs chan Job
	wg   sync.WaitGroup
}

func NewWorkerPool(workerCount int) *WorkerPool {
	wp := &WorkerPool{
		jobs: make(chan Job, workerCount),
	}

	for i := 0; i < workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}

	return wp
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	client := &http.Client{}

	for job := range wp.jobs {
		ctx, cancel := context.WithTimeout(job.Req.Context(), 5*time.Second)
		defer cancel()

		req := job.Req.WithContext(ctx)

		resp, err := client.Do(req)
		if err != nil {
			job.ErrorChan <- err
			continue
		}
		job.ResultChan <- resp
	}
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.jobs <- job
}

func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
}
