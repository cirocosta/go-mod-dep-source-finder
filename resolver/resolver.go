package resolver

import (
	"context"
	"net/http"

	"github.com/cirocosta/go-mod-license-finder/parser"
	"github.com/pkg/errors"
)

type Location struct {
	VCS string
	URL string
}

func doRequest(ctx context.Context, url string) (resp *http.Response, err error) {
	var req *http.Request

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		err = errors.Wrapf(err, "couldn't create request obj")
		return
	}

	req = req.WithContext(ctx)
	queryParams := req.URL.Query()
	queryParams.Add("go-get", "1")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		err = errors.Wrapf(err, "failed to send http request")
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		err = errors.Errorf("non-success status code %d", resp.StatusCode)
		return
	}

	return
}

func Resolve(ctx context.Context, dependency parser.Line) (loc Location, err error) {
	_, err = doRequest(ctx, dependency.Dependency)
	if err != nil {
		err = errors.Wrapf(err, "failed to issue request for dependency - %+v", dependency)
		return
	}

	return
}
