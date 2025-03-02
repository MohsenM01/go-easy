package httpclient

import (
    "context"
    "errors"
    "net/http"
    "time"
)

type HTTPClient struct {
    client      *http.Client
    maxRetries  int
    retryDelay  time.Duration
}

func NewHTTPClient(timeout time.Duration, maxRetries int, retryDelay time.Duration) *HTTPClient {
    return &HTTPClient{
        client: &http.Client{Timeout: timeout},
        maxRetries: maxRetries,
        retryDelay: retryDelay,
    }
}

func (hc *HTTPClient) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
    var resp *http.Response
    var err error

    for i := 0; i < hc.maxRetries; i++ {
        reqWithCtx := req.Clone(ctx)

        resp, err = hc.client.Do(reqWithCtx)
        if err == nil && resp.StatusCode < 500 { 
            return resp, nil
        }

        select {
        case <-ctx.Done(): 
            return nil, ctx.Err()
        case <-time.After(hc.retryDelay):
        }
    }

    return nil, errors.New("request failed after retries")
}
