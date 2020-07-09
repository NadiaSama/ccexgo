package request

import (
	"context"
	"net/http"
)

//Do issue http request and calls f with the response. if ctx.Done is called
//during the request handling. Do cancels the request and wait s for f quit
//return ctx.Err
func Do(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	req = req.WithContext(ctx)
	return DoReqWithCtx(req, f)
}

//DoReqWithCtx like Do except the ctx extract from req.Context()
func DoReqWithCtx(req *http.Request, f func(*http.Response, error) error) error {
	ctx := req.Context()
	c := make(chan error, 1)
	go func() { c <- f(http.DefaultClient.Do(req)) }()
	select {
	case <-ctx.Done():
		<-c
		return ctx.Err()

	case err := <-c:
		return err
	}
}
