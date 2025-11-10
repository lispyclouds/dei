package pw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdenticonOf(t *testing.T) {
	identicon, err := identiconOf("test", "test")
	assert.NoError(t, err)

	assert.Equal(t, "╔░╝☂", identicon)
}
