package s3

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"orbit/internal/security"
)

var sseMagic = []byte("ORBSSE1\x00")

const sseMasterLabel = "orbit/s3/sse-kms/v1"

// encryptSSE performs envelope encryption:
// generates a random 256-bit DEK, encrypts payload with DEK (AES-256-GCM),
// wraps DEK with KEK derived from Orbit's MasterKey, and packs:
// [8B magic][12B wrapNonce][48B wrappedDEK (32B+16B tag)][12B dataNonce][ciphertext+tag]
func (s *S3Engine) encryptSSE(plain []byte) ([]byte, error) {
	if len(s.cfg.MasterKey) == 0 {
		return nil, errors.New("master key unavailable for SSE")
	}
	kekBlock, err := aes.NewCipher(security.DeriveKey(s.cfg.MasterKey, sseMasterLabel))
	if err != nil {
		return nil, err
	}
	kekAEAD, err := cipher.NewGCM(kekBlock)
	if err != nil {
		return nil, err
	}

	dek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, dek); err != nil {
		return nil, err
	}
	wrapNonce := make([]byte, kekAEAD.NonceSize())
	if _, err := io.ReadFull(rand.Reader, wrapNonce); err != nil {
		return nil, err
	}
	wrappedDEK := kekAEAD.Seal(nil, wrapNonce, dek, nil)

	dekBlock, err := aes.NewCipher(dek)
	if err != nil {
		return nil, err
	}
	dekAEAD, err := cipher.NewGCM(dekBlock)
	if err != nil {
		return nil, err
	}
	dataNonce := make([]byte, dekAEAD.NonceSize())
	if _, err := io.ReadFull(rand.Reader, dataNonce); err != nil {
		return nil, err
	}
	ciphertext := dekAEAD.Seal(nil, dataNonce, plain, nil)

	out := make([]byte, 0, len(sseMagic)+len(wrapNonce)+len(wrappedDEK)+len(dataNonce)+len(ciphertext))
	out = append(out, sseMagic...)
	out = append(out, wrapNonce...)
	out = append(out, wrappedDEK...)
	out = append(out, dataNonce...)
	out = append(out, ciphertext...)
	return out, nil
}

// isSSEEncrypted checks whether data starts with the Orbit SSE envelope header.
func isSSEEncrypted(data []byte) bool {
	return len(data) >= len(sseMagic) && bytes.Equal(data[:len(sseMagic)], sseMagic)
}

func (s *S3Engine) decryptSSE(packed []byte) ([]byte, error) {
	if !isSSEEncrypted(packed) {
		return packed, nil
	}
	if len(s.cfg.MasterKey) == 0 {
		return nil, errors.New("master key unavailable for SSE decryption")
	}
	kekBlock, err := aes.NewCipher(security.DeriveKey(s.cfg.MasterKey, sseMasterLabel))
	if err != nil {
		return nil, err
	}
	kekAEAD, err := cipher.NewGCM(kekBlock)
	if err != nil {
		return nil, err
	}

	off := len(sseMagic)
	wrapNonceLen := kekAEAD.NonceSize()
	wrappedDEKLen := 32 + kekAEAD.Overhead()
	if len(packed) < off+wrapNonceLen+wrappedDEKLen+12 {
		return nil, errors.New("corrupt SSE envelope")
	}

	wrapNonce := packed[off : off+wrapNonceLen]
	off += wrapNonceLen
	wrappedDEK := packed[off : off+wrappedDEKLen]
	off += wrappedDEKLen

	dek, err := kekAEAD.Open(nil, wrapNonce, wrappedDEK, nil)
	if err != nil {
		return nil, errors.New("SSE DEK unwrap failed")
	}

	dekBlock, err := aes.NewCipher(dek)
	if err != nil {
		return nil, err
	}
	dekAEAD, err := cipher.NewGCM(dekBlock)
	if err != nil {
		return nil, err
	}
	dataNonceLen := dekAEAD.NonceSize()
	if len(packed) < off+dataNonceLen {
		return nil, errors.New("corrupt SSE payload")
	}
	dataNonce := packed[off : off+dataNonceLen]
	ciphertext := packed[off+dataNonceLen:]
	return dekAEAD.Open(nil, dataNonce, ciphertext, nil)
}
