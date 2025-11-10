package pw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainKey(t *testing.T) {
	key, err := mainKey("test", "pass", SiteVariant("password"))
	assert.NoError(t, err)

	assert.Equal(t, key, []byte{
		51, 253, 82, 252, 68, 97, 191, 162, 127, 73, 153, 160, 52, 128, 204, 4, 183,
		190, 106, 180, 68, 126, 100, 94, 132, 141, 99, 143, 106, 211, 94, 245, 245,
		255, 195, 72, 28, 128, 197, 51, 99, 27, 125, 24, 54, 193, 223, 230, 118,
		181, 225, 236, 171, 104, 9, 158, 214, 23, 166, 89, 36, 174, 64, 112,
	})
}

func TestPassword(t *testing.T) {
	key, err := mainKey("test", "pass", SiteVariant("password"))
	assert.NoError(t, err)

	password, err := password(key, "site", 1, PASSWORD, MAXIMUM)
	assert.NoError(t, err)

	assert.Equal(t, password, "QsKBWAYdT9dh^AOGVA0.")
}
