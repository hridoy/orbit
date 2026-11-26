package security

import "testing"

func TestPasswordHashRoundTrip(t *testing.T) {
	h, err := HashPassword("correct-horse")
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyPassword("correct-horse", h) {
		t.Fatal("right password rejected")
	}
	if VerifyPassword("correct-hors", h) {
		t.Fatal("wrong password accepted")
	}
	if VerifyPassword("correct-horse", "garbage") {
		t.Fatal("malformed hash accepted")
	}
	if _, err := HashPassword("short"); err == nil {
		t.Fatal("short password accepted")
	}
}
