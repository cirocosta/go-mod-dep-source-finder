package parser

import (
	"os"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

// Line represents the contents of a line without
// any parsing applied to the fields beside filtering out
// unwanted content.
//
type Line struct {
	Dependency                  string
	Reference                   string
	RelativeDirectoryDependency string
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
		return
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

// isAboveV1 verifies whether the version supplied corresponds to
// a version that is >= v2.
//
// ps.: panics if `content` is not a valid semver.
//
func isAboveV1(content string) bool {
	ver, err := version.NewVersion(content)
	if err != nil {
		panic(errors.Wrapf(err, "must've provided a valid semver"))
	}

	segments := ver.Segments()
	return segments[0] > 1
}

// isNotIncompatible checks whether the version is considered incompatible
// or not.
func isNotIncompatible(content string) bool {
	ver, err := version.NewVersion(content)
	if err != nil {
		panic(errors.Wrapf(err, "must've provided a valid semver"))
	}

	return ver.Metadata() != "incompatible"
}

// isVersionedDependency checks if the dependency provided at `content` is
// a dependency that is supposed to have a version in the end.
//
// Ref: https://github.com/golang/go/wiki/Modules#semantic-import-versioning
//
func isVersionedDependency(content string) bool {
	return !strings.HasPrefix(content, "gopkg.in")
}

// StripVersionFromDependency removes the version that might exist as the last
// fields in the `require` path for such dependency.
//
// For instance, having in `go.mod`:
//
// 	require github.com/my/mod/v2 v2.0.0
//
// This method would take:
//
// 	github.com/my/mod/v2
//
// And return
//
// 	github.com/my/mod
//
// If the last piece is a proper semver version, erroring otherwise.
//
func StripVersionFromDependency(content string) (res string, err error) {
	fields := strings.Split(content, "/")
	lastField := fields[len(fields)-1]

	_, err = version.NewVersion(lastField)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to parse semver")
		return
	}

	res = strings.Join(fields[0:len(fields)-1], "/")
	return
}

func sanitizeGithubDependency(content string) string {
	fields := strings.Split(content, "/")
	return fields[0] + "/" + fields[1] + "/" + fields[2]
}

func SanitizeDependency(content string) string {
	if strings.HasPrefix(content, "github.com") {
		return sanitizeGithubDependency(content)
	}

	return content
}

// ParseLine parses a single dependency line, returning the struct
// that represents its contents, with the version already interpreted
// as a reference.
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
// In case of a replace directive being used, only the target is taken into
// considration.
//
//
func ParseLine(content string) (line Line, err error) {
	content = sanitizeReplaceDirective(content)
	fields := strings.Fields(content)

	if len(fields) == 0 {
		err = errors.Errorf("not enough fields in line '%s'", content)
		return
	}

	if len(fields) == 1 {
		_, err = os.Stat(fields[0])
		if os.IsNotExist(err) {
			err = errors.Wrapf(err, "stat '%s'", fields[0])
			return
		}

		line.RelativeDirectoryDependency = fields[0]
		return
	}

	dependency := fields[0]
	version := fields[1]

	line.Reference, err = ParseVersion(version)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to parse version in line")
		return
	}

	if isAboveV1(version) && isNotIncompatible(version) && isVersionedDependency(dependency) {
		dependency, err = StripVersionFromDependency(dependency)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to strip version from dependency")
			return
		}
	}

	line.Dependency = SanitizeDependency(dependency)

	return
}

func sanitizeReplaceDirective(content string) string {
	const separator = "=>"

	if !strings.Contains(content, separator) {
		return content
	}

	return strings.Split(content, separator)[1]
}
