package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func Encrypt(p, k string) (string, error) {
	block, err := aes.NewCipher([]byte(k))
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(p))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(p))

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func Decrypt(c, k string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(c)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(k))
	if err != nil {
		return "", err
	}

	if len(data) < aes.BlockSize {
		return "", errors.New("cipher text too short")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return string(data), nil
}
