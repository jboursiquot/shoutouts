package shoutouts_test

import (
	"testing"

	"github.com/jboursiquot/shoutouts"
	"github.com/stretchr/testify/assert"
)

func TestShoutout(t *testing.T) {
	assert.NotEmpty(t, shoutouts.New().ID)
}
