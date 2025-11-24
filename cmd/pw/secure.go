package pw

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func padMeMaybe(toPad []byte) ([]byte, error) {
	blockSize := 32
	length := len(toPad)

	if length > blockSize {
		return nil, fmt.Errorf("Can only handle keys of max 32 bytes, got %d", len(toPad))
	}

	if length == blockSize {
		return toPad, nil
	}

	paddingSize := blockSize - len(toPad)%blockSize
	padding := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)

	return append(toPad, padding...), nil
}

func makeCipher(key []byte) (cipher.AEAD, error) {
	key, err := padMeMaybe(key)
	if err != nil {
		return nil, err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	return gcm, nil
}

func encrypt(data, key []byte) ([]byte, error) {
	c, err := makeCipher(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, c.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return c.Seal(nonce, nonce, data, nil), nil
}

func decrypt(data, key []byte) ([]byte, error) {
	c, err := makeCipher(key)
	if err != nil {
		return nil, err
	}

	nonceSize := c.NonceSize()
	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	unsealed, err := c.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return unsealed, nil
}
