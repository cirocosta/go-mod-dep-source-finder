package parser

import (
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

// Line represents the contents of a line without
// any parsing applied to the fields beside filtering out
// unwanted content.
//
type Line struct {
	Dependency string
	Reference  string
}

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

// ParseLine parses a single dependency line, returning the struct
// that represents its contents, with the version already interpreted
// as a reference.
//
func ParseLine(content string) (line Line, err error) {
	if content == "" {
		err = errors.Errorf("line must not be empty")
		return
	}

	fields := strings.Fields(content)
	if len(fields) < 2 {
		err = errors.Errorf("not enough fields")
		return
	}

	line = Line{
		Dependency: fields[0],
	}

	line.Reference, err = ParseVersion(fields[1])
	if err != nil {
		err = errors.Wrapf(err,
			"failed to parse version in line")
		return
	}

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
