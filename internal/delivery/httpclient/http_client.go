package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	workerpool "go-easy/pkg/utils"
)

type HTTPClient struct {
	pool       *workerpool.WorkerPool
	maxRetries int
	retryDelay time.Duration
}

func NewHTTPClient(workerCount int, maxRetries int, retryDelay time.Duration) *HTTPClient {
	pool := workerpool.NewWorkerPool(workerCount)
	return &HTTPClient{
		pool:       pool,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

func (hc *HTTPClient) sendRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	resultChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	hc.pool.AddJob(workerpool.Job{
		Req:        req,
		ResultChan: resultChan,
		ErrorChan:  errorChan,
	})

	select {
	case resp := <-resultChan:
		return resp, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (hc *HTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	return hc.sendRequest(ctx, http.MethodGet, url, nil)
}

func (hc *HTTPClient) Post(ctx context.Context, url string, body []byte) (*http.Response, error) {
	return hc.sendRequest(ctx, http.MethodPost, url, bytes.NewBuffer(body))
}

func (hc *HTTPClient) Put(ctx context.Context, url string, body []byte) (*http.Response, error) {
	return hc.sendRequest(ctx, http.MethodPut, url, bytes.NewBuffer(body))
}

func (hc *HTTPClient) Delete(ctx context.Context, url string) (*http.Response, error) {
	return hc.sendRequest(ctx, http.MethodDelete, url, nil)
}

func (hc *HTTPClient) Patch(ctx context.Context, url string, body []byte) (*http.Response, error) {
	return hc.sendRequest(ctx, http.MethodPatch, url, bytes.NewBuffer(body))
}

func (hc *HTTPClient) Close() {
	hc.pool.Close()
}
