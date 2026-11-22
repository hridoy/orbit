package security

import (
	"crypto/rand"
	"testing"
	"time"
)

func TestVaultEncryptionCycle(t *testing.T) {
	key := make([]byte, 32)
	_, _ = rand.Read(key)

	vault, err := NewVault(key)
	if err != nil {
		t.Fatalf("NewVault error: %v", err)
	}

	rawText := "SOVEREIGN_ORBIT_SECRET_PAYLOAD_12345"
	cipherB64, nonceB64, err := vault.Encrypt([]byte(rawText))
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}

	decrypted, err := vault.Decrypt(cipherB64, nonceB64)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}

	if string(decrypted) != rawText {
		t.Fatalf("expected '%s', got '%s'", rawText, string(decrypted))
	}
}

func TestTokenManager(t *testing.T) {
	seed := make([]byte, 32)
	_, _ = rand.Read(seed)

	tm, err := NewTokenManager(seed)
	if err != nil {
		t.Fatalf("NewTokenManager error: %v", err)
	}

	// 1. Master key validation
	claims, err := tm.VerifyToken(tm.MasterKeyHex())
	if err != nil {
		t.Fatalf("Master key verify failed: %v", err)
	}
	if claims.Role != "admin" {
		t.Errorf("expected master key role admin, got %s", claims.Role)
	}

	// 2. Token issue & verify
	token, err := tm.IssueToken("tok_123", "operator", 1*time.Hour)
	if err != nil {
		t.Fatalf("IssueToken error: %v", err)
	}

	verified, err := tm.VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken error: %v", err)
	}
	if verified.TokenID != "tok_123" || verified.Role != "operator" {
		t.Errorf("claims mismatch: %+v", verified)
	}

	// 3. Expired token
	expiredTok, err := tm.IssueToken("tok_exp", "deployer", -1*time.Second)
	if err != nil {
		t.Fatalf("IssueToken error: %v", err)
	}
	_, err = tm.VerifyToken(expiredTok)
	if err != ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}

	// 4. Tampered token
	tampered := token + "tampered"
	_, err = tm.VerifyToken(tampered)
	if err != ErrInvalidToken {
		t.Fatalf("expected ErrInvalidToken for tampered token, got %v", err)
	}
}
