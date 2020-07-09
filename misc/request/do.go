package request

import (
	"context"
	"net/http"
)

//Do issue http request and calls f with the response. if ctx.Done is called
//during the request handling. Do cancels the request and wait s for f quit
//return ctx.Err
func Do(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	c := make(chan error, 1)
	req = req.WithContext(ctx)
	go func() { c <- f(http.DefaultClient.Do(req)) }()
	select {
	case <-ctx.Done():
		<-c
		return ctx.Err()

	case err := <-c:
		return err
	}
}
