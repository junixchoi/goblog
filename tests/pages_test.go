package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHomePage(t *testing.T) {
	baseURL := "http://localhost:3000"

	var (
		resp *http.Response
		err  error
	)
	resp, err = http.Get(baseURL + "/")

	assert.NoError(t, err, "err is empty")
	assert.Equal(t, 200, resp.StatusCode, "must return 200 status code")
}
