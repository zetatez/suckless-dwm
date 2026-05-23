package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptBytesAESGCM encrypts plaintext with AES-GCM and returns nonce|ciphertext.
// aad can be nil; if provided, the same aad must be used during decryption.
func EncryptBytesAESGCM(key, plaintext, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ct := gcm.Seal(nil, nonce, plaintext, aad)
	out := append(nonce, ct...)
	return out, nil
}

// DecryptBytesAESGCM decrypts nonce|ciphertext produced by EncryptBytesAESGCM.
func DecryptBytesAESGCM(key, data, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce := data[:nonceSize]
	ct := data[nonceSize:]
	pt, err := gcm.Open(nil, nonce, ct, aad)
	if err != nil {
		return nil, err
	}
	return pt, nil
}

// EncryptStringAESGCM encrypts plaintext with AES-GCM and returns base64(nonce|ciphertext).
func EncryptStringAESGCM(key []byte, plaintext string) (string, error) {
	out, err := EncryptBytesAESGCM(key, []byte(plaintext), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(out), nil
}

// DecryptStringAESGCM decrypts base64(nonce|ciphertext) produced by EncryptStringAESGCM.
func DecryptStringAESGCM(key []byte, b64 string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	pt, err := DecryptBytesAESGCM(key, data, nil)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}
