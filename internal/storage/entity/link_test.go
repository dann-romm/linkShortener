package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateLink(t *testing.T) {
	assert.Equal(t, CreateLink("google.com"), CreateLink("google.com"))
	assert.NotEqual(t, CreateLink("google.com"), CreateLink("google.co"))
}
