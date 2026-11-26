package security

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid or tampered authentication token")
	ErrTokenExpired = errors.New("token has expired")
)

type TokenClaims struct {
	TokenID   string `json:"jti"`
	Role      string `json:"role"` // 'admin', 'operator', 'deployer'
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

type TokenManager struct {
	privKey ed25519.PrivateKey
	pubKey  ed25519.PublicKey
}

func NewTokenManager(seed []byte) (*TokenManager, error) {
	if len(seed) != 32 {
		return nil, errors.New("seed must be 32 bytes for Ed25519")
	}

	privKey := ed25519.NewKeyFromSeed(seed)
	pubKey := privKey.Public().(ed25519.PublicKey)

	return &TokenManager{
		privKey: privKey,
		pubKey:  pubKey,
	}, nil
}

// IssueToken mints a compact Ed25519-signed token: <base64(claims)>.<base64(signature)>.
func (tm *TokenManager) IssueToken(tokenID, role string, ttl time.Duration) (string, error) {
	now := time.Now().Unix()
	var exp int64
	if ttl != 0 {
		exp = now + int64(ttl.Seconds())
	}
	claims := TokenClaims{
		TokenID:   tokenID,
		Role:      role,
		IssuedAt:  now,
		ExpiresAt: exp,
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to encode claims: %w", err)
	}

	sig := ed25519.Sign(tm.privKey, payloadBytes)

	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadBytes)
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	return fmt.Sprintf("%s.%s", payloadB64, sigB64), nil
}

// VerifyToken validates the Ed25519 signature and expiration, returning claims.
func (tm *TokenManager) VerifyToken(tokenString string) (*TokenClaims, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 2 {
		return nil, ErrInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, ErrInvalidToken
	}

	sigBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !ed25519.Verify(tm.pubKey, payloadBytes, sigBytes) {
		return nil, ErrInvalidToken
	}

	var claims TokenClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt > 0 && time.Now().Unix() > claims.ExpiresAt {
		return nil, ErrTokenExpired
	}

	return &claims, nil
}
