package db

import (
	"testing"
)

func TestCreateURL(t *testing.T) {

	tests := []struct{
		BaseURL string
		Path string
		Expected string
	}{
		{
			Expected: "",
		},



		{
			BaseURL: "http:",
			Expected: "",
		},



		{
			BaseURL: "//example.com",
			Expected: "",
		},
		{
			BaseURL: "//example.com/",
			Expected: "",
		},
		{
			BaseURL: "//example.com//",
			Expected: "",
		},
		{
			BaseURL: "//example.com/api",
			Expected: "",
		},
		{
			BaseURL: "//example.com/api/",
			Expected: "",
		},
		{
			BaseURL: "//example.com//api/",
			Expected: "",
		},
		{
			BaseURL: "//example.com/api//",
			Expected: "",
		},
		{
			BaseURL: "//example.com//api//",
			Expected: "",
		},



		{
			BaseURL:  "http://example.com",
			Expected: "http://example.com/",
		},
		{
			BaseURL:  "http://example.com/",
			Expected: "http://example.com/",
		},
		{
			BaseURL : "http://example.com//",
			Expected: "http://example.com/",
		},
		{
			BaseURL : "http://example.com///",
			Expected: "http://example.com/",
		},
		{
			BaseURL : "http://example.com////",
			Expected: "http://example.com/",
		},



		{
			BaseURL:  "http://example.com/api",
			Expected: "http://example.com/api",
		},
		{
			BaseURL:  "http://example.com/api/",
			Expected: "http://example.com/api/",
		},
		{
			BaseURL : "http://example.com//api",
			Expected: "http://example.com/api",
		},
		{
			BaseURL:  "http://example.com/api//",
			Expected: "http://example.com/api/",
		},
		{
			BaseURL:  "http://example.com//api//",
			Expected: "http://example.com/api/",
		},
		{
			BaseURL:  "http://example.com///api////",
			Expected: "http://example.com/api/",
		},



		{
			BaseURL:  "http://example.com",
			Path:                        "apple/banana/cherry",
			Expected: "http://example.com/apple/banana/cherry",
		},
		{
			BaseURL:  "http://example.com/",
			Path:                        "apple/banana/cherry",
			Expected: "http://example.com/apple/banana/cherry",
		},
		{
			BaseURL : "http://example.com//",
			Path:                        "apple/banana/cherry",
			Expected: "http://example.com/apple/banana/cherry",
		},



		{
			BaseURL:  "http://example.com",
			Path:                       "/apple/banana/cherry",
			Expected: "http://example.com/apple/banana/cherry",
		},
		{
			BaseURL:  "http://example.com/",
			Path:                       "/apple/banana/cherry",
			Expected: "http://example.com/apple/banana/cherry",
		},
		{
			BaseURL : "http://example.com//",
			Path:                       "/apple/banana/cherry",
			Expected: "http://example.com/apple/banana/cherry",
		},



		{
			BaseURL:  "http://example.com",
			Path:                       "/apple//banana/cherry",
			Expected: "http://example.com/apple/banana/cherry",
		},
		{
			BaseURL:  "http://example.com/",
			Path:                       "/apple//banana/cherry",
			Expected: "http://example.com/apple/banana/cherry",
		},
		{
			BaseURL : "http://example.com//",
			Path:                       "/apple//banana/cherry",
			Expected: "http://example.com/apple/banana/cherry",
		},



		{
			BaseURL:  "http://example.com",
			Path:                      "//apple//banana/cherry",
			Expected: "http://example.com/apple/banana/cherry",
		},
		{
			BaseURL:  "http://example.com/",
			Path:                      "//apple//banana/cherry",
			Expected: "http://example.com/apple/banana/cherry",
		},
		{
			BaseURL : "http://example.com//",
			Path:                      "//apple//banana/cherry",
			Expected: "http://example.com/apple/banana/cherry",
		},



		{
			BaseURL:  "http://example.com/api",
			Path:                            "apple/banana/cherry",
			Expected: "http://example.com/api/apple/banana/cherry",
		},
		{
			BaseURL:  "http://example.com/api/",
			Path:                            "apple/banana/cherry",
			Expected: "http://example.com/api/apple/banana/cherry",
		},
		{
			BaseURL : "http://example.com/api//",
			Path:                            "apple/banana/cherry",
			Expected: "http://example.com/api/apple/banana/cherry",
		},
		{
			BaseURL : "http://example.com//api//",
			Path:                            "apple/banana/cherry",
			Expected: "http://example.com/api/apple/banana/cherry",
		},
	}

	for testNumber, test := range tests {

		actual := createURL(test.BaseURL, test.Path)

		expected := test.Expected

		if expected != actual {
			t.Errorf("For test #%d, the actual created-URL is not what was expected.", testNumber)
			t.Logf("EXPECTED: %q", expected)
			t.Logf("ACTUAL:   %q", actual)
			t.Logf("BASE-URL: %q", test.BaseURL)
			t.Logf("PATH: %q", test.Path)
			continue
		}
	}
}

