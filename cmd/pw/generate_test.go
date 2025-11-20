package pw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnlyHosts(t *testing.T) {
	assert.Equal(t, "youtube.com", onlyHosts("https://www.youtube.com/watch?v=t6qL_FbLArk"))
	assert.Equal(t, "github.com", onlyHosts("https://github.com/babashka/pod-babashka-go-sqlite3/blob/main/CHANGELOG.md"))
	assert.Equal(t, "github.com", onlyHosts("github.com"))
}
