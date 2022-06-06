package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestAllPage(t *testing.T) {
	baseURL := "http://localhost:3000"

	// init test data
	var tests = []struct {
		method   string
		url      string
		expected int
	}{
		{"GET", "/", 200},
		{"GET", "/about", 200},
		{"GET", "/notfound", 404},
		{"GET", "/articles", 200},
		{"GET", "/articles/create", 200},
		{"GET", "/articles/3", 200},
		{"GET", "/articles/3/edit", 200},
		{"POST", "/articles/3", 200},
		{"POST", "/articles", 200},
		{"POST", "/articles/1/delete", 404},
	}

	// 2. loop all tests
	for _, test := range tests {
		t.Logf("current request URL: %v \n", test.url)
		var (
			resp *http.Response
			err  error
		)
		// 2.1 request and get response
		switch {
		case test.method == "POST":
			data := make(map[string][]string)
			resp, err = http.PostForm(baseURL+test.url, data)
		default:
			resp, err = http.Get(baseURL + test.url)
		}
		assert.NoError(t, err, "request "+test.url+" error")
		assert.Equal(t, test.expected, resp.StatusCode, test.url+" must return status code "+strconv.Itoa(test.expected))
	}
}
