package parser

import (
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

// ParseVersion retrieves a reference out of a version.
//
func ParseVersion(ver string) (ref string, err error) {
	if ver == "" {
		err = errors.Errorf("version must not be empty")
		return
	}

	_, err = version.NewVersion(ver)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to parse semver")
	}

	ref = strings.TrimSuffix(ver, "+incompatible")

	var (
		ok       bool
		obtained string
	)

	if obtained, ok = tryVersionDateSha(ref); ok {
		ref = obtained
		return
	}

	return
}

// tryVersionDateSha tries to parse the version as a pseudo-version
// following `vX.0.0-yyyymmddhhmmss-abcdefabcdef`.
//
func tryVersionDateSha(ver string) (ref string, ok bool) {
	components := strings.Split(ver, "-")

	if len(components) != 3 {
		return
	}

	ok = true
	ref = components[2]

	return
}

// Line represents the contents of a line without
// any parsing applied to the fields beside filtering out
// unwanted content.
//
type Line struct {
	Dependency string
	Version    string
}

func ParseLine(content string) (line Line, err error) {
	// TODO
	return
}

// Parse parses a go.mod dependency line.
//
//
// 	gopkg.in/gorethink/gorethink.v4 v4.1.0+incompatible // indirect
//      |                             | |    | |         |  |         |
//      *-+---------------------------* *-+--* *+--------*  *-----+---*
//        |                               |     |                 |
//        |                               |     excluded from parsing
//        |                               |
//        |                             reference (e.g., git tag)
//        |
//        dependency name (where to ask for go source code)
//
//
// func Parse() (dep Dependency, err error) {
// 	return
// }
