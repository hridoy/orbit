package security

import (
	"strings"
	"testing"
	"time"
)

// RFC 6238 Appendix B test secret ("12345678901234567890") in base32.
const rfcSecret = "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"

func TestTOTPRFC6238Vectors(t *testing.T) {
	// The RFC lists 8-digit SHA-1 codes; the 6-digit code is the last six digits.
	cases := []struct {
		unix int64
		want string
	}{
		{59, "287082"},
		{1111111109, "081804"},
		{1111111111, "050471"},
		{1234567890, "005924"},
		{2000000000, "279037"},
	}
	for _, c := range cases {
		got, err := TOTPCode(rfcSecret, time.Unix(c.unix, 0))
		if err != nil {
			t.Fatal(err)
		}
		if got != c.want {
			t.Fatalf("time %d: got %s want %s", c.unix, got, c.want)
		}
	}
}

func TestVerifyTOTPSkewAndReplayStep(t *testing.T) {
	now := time.Unix(1111111109, 0)
	code, _ := TOTPCode(rfcSecret, now)
	step, ok := VerifyTOTP(rfcSecret, code, now)
	if !ok || step != now.Unix()/30 {
		t.Fatalf("current code rejected: ok=%v step=%d", ok, step)
	}
	// One step earlier and later are accepted (clock drift).
	prev, _ := TOTPCode(rfcSecret, now.Add(-30*time.Second))
	if _, ok := VerifyTOTP(rfcSecret, prev, now); !ok {
		t.Fatal("previous step should be accepted")
	}
	// Two steps away is rejected.
	far, _ := TOTPCode(rfcSecret, now.Add(-90*time.Second))
	if _, ok := VerifyTOTP(rfcSecret, far, now); ok {
		t.Fatal("code two steps away must be rejected")
	}
}

func TestVerifyTOTPRejectsMalformed(t *testing.T) {
	now := time.Unix(59, 0)
	for _, bad := range []string{"", "12345", "1234567", "abcdef"} {
		if _, ok := VerifyTOTP(rfcSecret, bad, now); ok {
			t.Fatalf("accepted malformed code %q", bad)
		}
	}
}

func TestGenerateTOTPSecretAndURI(t *testing.T) {
	s, err := GenerateTOTPSecret()
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 32 {
		t.Fatalf("secret length %d, want 32 base32 chars", len(s))
	}
	uri := TOTPProvisioningURI(s, "ravi")
	if !strings.HasPrefix(uri, "otpauth://totp/Orbit:ravi?") || !strings.Contains(uri, "secret="+s) {
		t.Fatalf("bad URI: %s", uri)
	}
}
