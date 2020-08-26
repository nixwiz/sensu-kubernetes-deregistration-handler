package main

import (
	"html"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
}

func TestCheckArgs(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	assert.Error(checkArgs(event))
	plugin.APIURL = "NotURL"
	assert.Error(checkArgs(event))
	plugin.APIURL = "http://127.0.0.1:8080"
	assert.Error(checkArgs(event))
	plugin.APIKey = "01a23bc4-56d7-890e-fa12-3456789bcd01"
	assert.NoError(checkArgs(event))
}

func TestExecuteHandler(t *testing.T) {
	testcases := []struct {
		checkName   string
		expectError bool
		httpStatus int
	}{
		{"check1", true, 200},
		{"kubernetes-delete-entity", false, 200},
		{"kubernetes-delete-entity", false, 404},
	}
	for _, tc := range testcases {
		assert := assert.New(t)
		event := corev2.FixtureEvent("entity1", tc.checkName)
		plugin.APIKey = "blah-blah-blah"

		var test = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(html.EscapeString(r.URL.Path), "/api/core/v2/namespaces/default/entities/entity1")
			_, err := ioutil.ReadAll(r.Body)
			assert.NoError(err)
			w.WriteHeader(tc.httpStatus)
		}))
		_, err := url.ParseRequestURI(test.URL)
		require.NoError(t, err)
		plugin.APIURL = test.URL
		if tc.expectError {
			assert.Error(executeHandler(event))
		} else {
			assert.NoError(executeHandler(event))
		}
	}
}
