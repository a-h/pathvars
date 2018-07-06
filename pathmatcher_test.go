package pathvars

import (
	"net/url"
	"reflect"
	"testing"
)

func TestPathExtraction(t *testing.T) {
	tests := []struct {
		name              string
		patterns          []string
		inputURL          string
		expectedVariables map[string]string
		expectedMatch     bool
	}{
		{
			name: "match single variable, no trailing slash",
			patterns: []string{
				"/user/{userid}",
			},
			inputURL: "/user/123",
			expectedVariables: map[string]string{
				"userid": "123",
			},
			expectedMatch: true,
		},
		{
			name: "full URL, not just the path",
			patterns: []string{
				"/user/{userid}",
			},
			inputURL: "https://subdomain.example.com/user/123",
			expectedVariables: map[string]string{
				"userid": "123",
			},
			expectedMatch: true,
		},
		{
			name: "match single variable, with trailing slash",
			patterns: []string{
				"/user/{userid}",
			},
			inputURL: "/user/123/",
			expectedVariables: map[string]string{
				"userid": "123",
			},
			expectedMatch: true,
		},
		{
			name: "match multiple variables",
			patterns: []string{
				"/organisation/{orgid}/user/{userid}/edit",
			},
			inputURL: "/organisation/123/user/456/edit",
			expectedVariables: map[string]string{
				"orgid":  "123",
				"userid": "456",
			},
			expectedMatch: true,
		},
		{
			name: "multiple patterns, no match",
			patterns: []string{
				"/user/{userid}",
				"/another/url",
				"/another/url/with/multiple/",
			},
			inputURL:          "/something/123/",
			expectedVariables: nil,
			expectedMatch:     false,
		},
		{
			name: "wildcard prefix",
			patterns: []string{
				"*/{userid}",
			},
			inputURL: "/something/123/",
			expectedVariables: map[string]string{
				"userid": "123",
			},
			expectedMatch: true,
		},
		{
			name: "multi-level wildcard prefix 1",
			patterns: []string{
				"*/{userid}",
			},
			inputURL: "/prefix/something/123/",
			expectedVariables: map[string]string{
				"userid": "123",
			},
			expectedMatch: true,
		},
		{
			name: "multi-level wildcard prefix 2",
			patterns: []string{
				"*/thing/{userid}",
			},
			inputURL: "/another/prefix/something/thing/123/",
			expectedVariables: map[string]string{
				"userid": "123",
			},
			expectedMatch: true,
		},
		{
			name: "wildcard suffix",
			patterns: []string{
				"*/includes/{id}/*",
			},
			inputURL: "/something/prefix/includes/123/this/",
			expectedVariables: map[string]string{
				"id": "123",
			},
			expectedMatch: true,
		},
		{
			name: "mismatched wildcard",
			patterns: []string{
				"*/includes/{id}/*",
			},
			inputURL:      "/something/prefix/notincluded/123/this/",
			expectedMatch: false,
		},
	}

	for _, test := range tests {
		e := NewExtractor(test.patterns...)
		u, err := url.Parse(test.inputURL)
		if err != nil {
			t.Errorf("%s: failed to parse URL '%s' with error: %v", test.name, test.inputURL, err)
			continue
		}
		actualVariables, actualMatch := e.Extract(u)
		if actualMatch != test.expectedMatch {
			t.Errorf("%s: expected match %v, got %v", test.name, test.expectedMatch, actualMatch)
		}
		if !reflect.DeepEqual(actualVariables, test.expectedVariables) {
			t.Errorf("%s: expected variables %v, got %v", test.name, test.expectedVariables, actualVariables)
		}
	}
}
