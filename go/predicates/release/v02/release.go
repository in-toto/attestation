/*
Validator APIs for release v0.2 protos.
*/
package v02

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrPurlRequired        = errors.New("purl required")
	ErrInvalidPurl         = errors.New("purl must be a package URL of the form pkg:type/name@version")
	ErrPurlVersionRequired = errors.New("purl must include a version")
)

// Validate checks that purl is set and is a package URL with a version, as
// the release predicate's parsing rules require. Only the parts needed for
// that are parsed: the pkg scheme, a type and a name, and the version. The
// purl spec itself makes the version optional, the release predicate
// doesn't.
func (r *Release) Validate() error {
	purl := r.GetPurl()
	if purl == "" {
		return ErrPurlRequired
	}
	version, err := purlVersion(purl)
	if err != nil {
		return fmt.Errorf("%w, got %q", err, purl)
	}
	if version == "" {
		return fmt.Errorf("%w, got %q", ErrPurlVersionRequired, purl)
	}
	return nil
}

// purlVersion returns the version of a package URL, following the parsing
// order of the purl spec: the subpath after '#' and the qualifiers after '?'
// come off first, then the version after the last '@' of the name.
func purlVersion(purl string) (string, error) {
	scheme, rest, ok := strings.Cut(purl, ":")
	if !ok || !strings.EqualFold(scheme, "pkg") {
		return "", ErrInvalidPurl
	}
	rest, _, _ = strings.Cut(rest, "#")
	rest, _, _ = strings.Cut(rest, "?")
	rest = strings.Trim(rest, "/")

	// The version is only looked for in the last segment, so that an
	// unencoded '@' in a namespace, like npm's @angular, isn't taken for it.
	var version string
	last := strings.LastIndex(rest, "/") + 1
	if i := strings.LastIndex(rest[last:], "@"); i >= 0 {
		rest, version = rest[:last+i], rest[last+i+1:]
	}

	// Both the type and the name are required.
	typ, name, ok := strings.Cut(rest, "/")
	if !ok || typ == "" || strings.Trim(name, "/") == "" {
		return "", ErrInvalidPurl
	}
	return version, nil
}
