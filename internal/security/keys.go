package security

import (
	"crypto/hmac"
	"crypto/sha256"
)

// DeriveKey returns a 32-byte subkey bound to a purpose label. Subkeys are domain-separated,
// so compromising one purpose does not reveal the master key or other subkeys.
func DeriveKey(masterKey []byte, label string) []byte {
	mac := hmac.New(sha256.New, masterKey)
	mac.Write([]byte(label))
	return mac.Sum(nil)
}

// TokenSeedLabel is the derivation label for Ed25519 token signing.
const TokenSeedLabel = "orbit/token-seed/v1"
