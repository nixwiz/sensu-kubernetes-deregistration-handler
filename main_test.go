package main

import (
	"testing"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
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
