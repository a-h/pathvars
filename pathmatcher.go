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
	s := u.Path
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

// ExtractString extracts path variables from a string.
func (pm *Extractor) ExtractString(s string) (v map[string]string, ok bool) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}
	return pm.Extract(u)
}
