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
			name: "multi-level wildcard prefix 3",
			patterns: []string{
				"*/thing/{userid}",
			},
			inputURL: "/thing/123",
			expectedVariables: map[string]string{
				"userid": "123",
			},
			expectedMatch: true,
		},
		{
			name: "multi-level wildcard prefix 4",
			patterns: []string{
				"*/another/thing/{userid}",
			},
			inputURL:      "/thing/123",
			expectedMatch: false,
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
		{
			name: "path encoded values are correctly decoded",
			patterns: []string{
				"/file/{filename}",
			},
			inputURL: "/file/this%20is%20a%20file.txt",
			expectedVariables: map[string]string{
				"filename": "this is a file.txt",
			},
			expectedMatch: true,
		},
		{
			name: "slashes in path encoded values are correctly decoded",
			patterns: []string{
				"/file/{filename}",
			},
			inputURL: "/file/this%2Fis%2Fa%2Ffile.txt",
			expectedVariables: map[string]string{
				"filename": "this/is/a/file.txt",
			},
			expectedMatch: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := NewExtractor(test.patterns...)
			actualVariables, actualMatch := e.ExtractString(test.inputURL)
			if actualMatch != test.expectedMatch {
				t.Errorf("%s: expected match %v, got %v", test.name, test.expectedMatch, actualMatch)
			}
			if !reflect.DeepEqual(actualVariables, test.expectedVariables) {
				t.Errorf("%s: expected variables %v, got %v", test.name, test.expectedVariables, actualVariables)
			}
		})
	}
}

func BenchmarkPathExtraction(b *testing.B) {
	e := NewExtractor("/user/{userid}")
	inputURL, err := url.Parse("/user/123")
	if err != nil {
		b.Fatalf("failed to parse URL: %v", err)
	}
	expectedVariables := map[string]string{
		"userid": "123",
	}
	expectedMatch := true

	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		actualVariables, actualMatch := e.Extract(inputURL)
		if actualMatch != expectedMatch {
			b.Fatalf("expected match %v, got %v", expectedMatch, actualMatch)
		}
		if !reflect.DeepEqual(actualVariables, expectedVariables) {
			b.Fatalf("expected variables %v, got %v", expectedVariables, actualVariables)
		}
	}
}
