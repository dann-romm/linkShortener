package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateLink(t *testing.T) {
	assert.Equal(t, CreateLink("google.com"), CreateLink("google.com"))
	assert.NotEqual(t, CreateLink("google.com"), CreateLink("google.co"))
	assert.Equal(t, CreateLink(""), CreateLink(""))
}

func TestNewLink(t *testing.T) {
	_, err := NewLink("")
	assert.Error(t, err)

	link1, err := NewLink("google.com")
	assert.NoError(t, err)
	link2, err := NewLink("google.com")
	assert.NoError(t, err)
	assert.Equal(t, link1, link2)
}

func TestValidate(t *testing.T) {
	assert.True(t, validate(0xFBDFEC0BFD7EAF9F))
	assert.True(t, validate(0x7FE))
	assert.False(t, validate(0xFFE))
	assert.False(t, validate(0x3F))
	assert.True(t, validate(0x0))

}
