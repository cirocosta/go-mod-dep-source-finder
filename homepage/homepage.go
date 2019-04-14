package homepage

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func Find(address, version string) (homepage string, unknownHost bool, err error) {
	parsedUrl, err := url.Parse(address)
	if err != nil {
		err = errors.Wrapf(err,
			"failed parsing url")
		return
	}

	switch parsedUrl.Host {
	case "github.com", "gitlab.com":
		homepage = address + "/tree/" + version
		return
	case "bitbucket.org":
		homepage = address + "/src/" + version
		return
	case "go.googlesource.com", "code.googlesource.com":
		homepage = address + "/+/" + version + "/"
		return
	case "gopkg.in":
		fields := strings.Split(strings.Trim(parsedUrl.Path, "/"), "/")
		lastField := fields[len(fields)-1]

		lastFieldFields := strings.Split(lastField, ".")
		packageName := lastFieldFields[0]

		if len(fields) == 1 {
			homepage = "https://github.com/go-" + packageName + "/" + packageName
		} else {
			homepage = "https://github.com/" + fields[0] + "/" + packageName
		}

		homepage = homepage + "/tree/" + version
		return
	}

	unknownHost = true
	return
}
