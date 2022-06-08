package linkservice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// test transformLink function
func TestTransformLink(t *testing.T) {
	assert.Equal(t, transformLink("google.com"), "google.com")
	assert.Equal(t, transformLink("google.co"), "google.co")
	assert.Equal(t, transformLink("google.com/"), "google.com")
	assert.Equal(t, transformLink("http://google.com"), "google.com")
	assert.Equal(t, transformLink("https://google.com"), "google.com")
	assert.Equal(t, transformLink("https://google.com/////"), "google.com")
	assert.Equal(t, transformLink("https://google.co/m/////"), "google.co/m")
	assert.Equal(t, transformLink("///https://google.co/m"), "///https://google.co/m")
	assert.Equal(t, transformLink("/////"), "")
	assert.Equal(t, transformLink("https://///"), "")
	assert.Equal(t, transformLink("https:/google.com/"), "https:/google.com")
}
