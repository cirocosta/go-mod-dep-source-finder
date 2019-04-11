package resolver

import (
	"context"

	"github.com/cirocosta/go-mod-license-finder/parser"
)

type Location struct {
	VCS string
	URL string
}

func Resolve(ctx context.Context, dependency parser.Line) (loc Location, err error) {
	return
}
