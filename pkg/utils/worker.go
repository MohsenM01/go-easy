package workerpool

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	httpclient "go-easy/internal/delivery/http"
)

type Job struct {
	URL    string
	Method string
	Body   []byte
}

type WorkerPool struct {
	client *httpclient.HTTPClient
	jobs   chan Job
	wg     sync.WaitGroup
}

func NewWorkerPool(workerCount int, client *httpclient.HTTPClient) *WorkerPool {
	wp := &WorkerPool{
		client: client,
		jobs:   make(chan Job, workerCount),
	}

	for i := 0; i < workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}

	return wp
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for job := range wp.jobs {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // محدودیت زمانی برای هر درخواست
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, job.Method, job.URL, bytes.NewBuffer(job.Body))
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := wp.client.DoRequest(ctx, req)
		if err != nil {
			fmt.Println("Request failed:", err)
			continue
		}

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response:", string(body))
		resp.Body.Close()
	}
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.jobs <- job
}

func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
}
