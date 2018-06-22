package pathvars

import "strings"

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
	if len(segments) != len(r) {
		return
	}
	vars = make(map[string]string)
	for i, inputSegment := range segments {
		routeSegment := r[i]
		name, capture, matches := routeSegment.Match(inputSegment)
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
