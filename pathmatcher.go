package pathvars

import (
	"net/url"
	"strings"
)

// Extractor checks routes to see if any can be matched, and what path variables are in them.
type Extractor struct {
	Routes []*Route
}

// NewExtractor creates an extractor.
func NewExtractor(patterns ...string) *Extractor {
	e := &Extractor{
		Routes: make([]*Route, len(patterns)),
	}
	for i, p := range patterns {
		e.Routes[i] = NewRoute(p)
	}
	return e
}

// Extract variables from the path.
func (pm *Extractor) Extract(u *url.URL) (v map[string]string, ok bool) {
	s := u.String()
	s = strings.TrimSuffix(s, "/")
	s = strings.TrimPrefix(s, "/")
	segments := strings.Split(s, "/")

	for _, r := range pm.Routes {
		v, ok = r.Match(segments)
		if ok {
			return
		}
	}
	v = nil
	return
}
