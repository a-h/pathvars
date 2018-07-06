package pathvars

import (
	"strings"
)

// Route is an array of segments.
type Route []*Segment

// NewRoute creates a route based on a pattern, e.g /users/{userid}.
func NewRoute(pattern string) *Route {
	var r Route

	pattern = strings.TrimSuffix(pattern, "/")
	pattern = strings.TrimPrefix(pattern, "/")

	for _, seg := range strings.Split(pattern, "/") {
		ps := &Segment{
			Name: seg,
		}
		if seg == "*" {
			ps.IsWildcard = true
		}
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			ps.IsVariable = true
			ps.Name = strings.TrimSuffix(strings.TrimPrefix(seg, "{"), "}")
		}
		r = append(r, ps)
	}

	return &r
}

// Match returns whether the route was matched, and extracts variables.
func (r Route) Match(segments []string) (vars map[string]string, ok bool) {
	vars = make(map[string]string)
	var wildcard bool
	for i := 0; i < len(r); i++ {
		routeSegment := r[len(r)-1-i]
		inputSegment := segments[len(segments)-1-i]
		name, capture, wildcardMatch, matches := routeSegment.Match(inputSegment)
		if matches {
			if wildcardMatch {
				wildcard = true
			} else {
				wildcard = false
			}
		}
		if wildcard {
			matches = true
		}
		if !matches {
			return
		}
		if capture {
			vars[name] = inputSegment
		}
	}
	ok = true
	return
}
