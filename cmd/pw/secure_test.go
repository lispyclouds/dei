package pw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	plainText := "foo bar baz"
	key := "secureKey42069"

	encrypted, err := encrypt([]byte(plainText), []byte(key))
	assert.NoError(t, err)

	decrypted, err := decrypt(encrypted, []byte(key))
	assert.NoError(t, err)

	assert.Equal(t, plainText, string(decrypted))
}
