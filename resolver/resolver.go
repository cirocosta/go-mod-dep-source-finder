package resolver

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type Location struct {
	VCS string
	URL string
}

// doRequest performs an HTTP request against a URL with the necessary setup for
// capturing pages with `go-import`s.
//
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

// GoImport is a representation of the contents that get parsed from
// a `go-import` meta element in an HTML page.
//
// For instance:
//
// 	<meta
// 		name="go-import"
// 		content="gopkg.in/yaml.v2 git https://gopkg.in/yaml.v2" --> what GoImport represents
// 	>
//
//
type GoImport struct {
	ImportPrefix string
	VCS          string
	RepoRoot     string
}

// ParseGoImport parses the value associated with the content tag of a `go-import`
// header meta element.
//
func ParseGoImport(content string) (res GoImport, err error) {
	if content == "" {
		err = errors.Errorf("empty content")
		return
	}

	fields := strings.Fields(content)
	if len(fields) != 3 {
		err = errors.Errorf("must have 3 fields")
		return
	}

	res.ImportPrefix = fields[0]
	res.VCS = fields[1]
	res.RepoRoot = fields[2]

	return
}

// FindGoImport searches for a `go-import` through the contents of a reader, capturing
// the `go-import` content's if found.
//
// Ref: https://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies
//
func FindGoImport(reader io.Reader) (importLine string, found bool, err error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		err = errors.Wrapf(err, "failed to create document from reader")
		return
	}

	doc.Find(`meta[name="go-import"]`).Each(func(i int, s *goquery.Selection) {
		content, exists := s.Attr("content")
		if !exists {
			return
		}

		found = true
		importLine = content
	})

	return
}

// Resolve retrieves a dependency's source code location from a dependency line.
//
func Resolve(ctx context.Context, dependency string) (loc Location, err error) {
	var resp *http.Response

	resp, err = doRequest(ctx, dependency)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to issue request for dependency - %+v",
			dependency,
		)
		return
	}

	defer resp.Body.Close()

	goImportContent, found, err := FindGoImport(resp.Body)
	if err != nil {
		err = errors.Wrapf(err, "failed to find `go-import` in body")
		return
	}

	if !found {
		err = errors.Errorf("import line not found")
		return
	}

	goImport, err := ParseGoImport(goImportContent)
	if err != nil {
		err = errors.Wrapf(err, "failed parsing go import content")
		return
	}

	loc.URL = goImport.RepoRoot
	loc.VCS = goImport.VCS

	return
}
