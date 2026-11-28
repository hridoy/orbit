package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// RFC 6238 time-based one-time passwords as used by authenticator apps: HMAC-SHA1, 6 digits,
// 30-second steps. One step of clock drift is accepted in either direction.
const (
	totpStep   = 30
	totpDigits = 6
	totpSkew   = 1
	totpIssuer = "Orbit"
)

var totpEncoding = base32.StdEncoding.WithPadding(base32.NoPadding)

// GenerateTOTPSecret returns a new random 160-bit secret encoded as base32, the form authenticator
// apps expect.
func GenerateTOTPSecret() (string, error) {
	buf := make([]byte, 20)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return totpEncoding.EncodeToString(buf), nil
}

// TOTPCode returns the 6-digit code for the given base32 secret at a moment in time.
func TOTPCode(secret string, t time.Time) (string, error) {
	return totpCodeForStep(secret, t.Unix()/totpStep)
}

func totpCodeForStep(secret string, step int64) (string, error) {
	key, err := totpEncoding.DecodeString(strings.ToUpper(strings.TrimSpace(secret)))
	if err != nil {
		return "", fmt.Errorf("invalid TOTP secret: %w", err)
	}
	var msg [8]byte
	binary.BigEndian.PutUint64(msg[:], uint64(step))
	mac := hmac.New(sha1.New, key)
	mac.Write(msg[:])
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	code := binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7fffffff
	mod := uint32(1)
	for i := 0; i < totpDigits; i++ {
		mod *= 10
	}
	return fmt.Sprintf("%0*d", totpDigits, code%mod), nil
}

// VerifyTOTP checks a submitted code against the secret within the allowed clock skew. It returns
// the matched time step so the caller can refuse to accept the same step twice (replay).
func VerifyTOTP(secret, code string, now time.Time) (int64, bool) {
	code = strings.TrimSpace(code)
	if len(code) != totpDigits {
		return 0, false
	}
	current := now.Unix() / totpStep
	for d := -totpSkew; d <= totpSkew; d++ {
		step := current + int64(d)
		want, err := totpCodeForStep(secret, step)
		if err != nil {
			return 0, false
		}
		if hmac.Equal([]byte(want), []byte(code)) {
			return step, true
		}
	}
	return 0, false
}

// TOTPProvisioningURI builds the otpauth:// URI that authenticator apps read from a QR code or
// a pasted setup key.
func TOTPProvisioningURI(secret, account string) string {
	label := url.PathEscape(totpIssuer) + ":" + url.PathEscape(account)
	q := url.Values{}
	q.Set("secret", secret)
	q.Set("issuer", totpIssuer)
	q.Set("algorithm", "SHA1")
	q.Set("digits", fmt.Sprint(totpDigits))
	q.Set("period", fmt.Sprint(totpStep))
	return "otpauth://totp/" + label + "?" + q.Encode()
}
