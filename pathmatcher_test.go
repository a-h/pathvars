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
