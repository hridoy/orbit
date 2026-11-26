package tlsengine

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"orbit/internal/db"
	"orbit/internal/security"
)

func newEngine(t *testing.T) *Engine {
	t.Helper()
	conn, err := db.Open(filepath.Join(t.TempDir(), "orbit.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	if err := conn.Migrate(); err != nil {
		t.Fatal(err)
	}
	key := make([]byte, 32)
	key[0] = 7
	v, err := security.NewVault(key)
	if err != nil {
		t.Fatal(err)
	}
	return New(db.NewTLSDomainStore(conn), v)
}

// selfSigned makes a certificate and key for domain, as PEM.
func selfSigned(t *testing.T, domain string) (string, string) {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: domain},
		DNSNames:     []string{domain},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(90 * 24 * time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &k.PublicKey, k)
	if err != nil {
		t.Fatal(err)
	}
	keyDER, err := x509.MarshalECPrivateKey(k)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})),
		string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}))
}

func TestNormalizeDomain(t *testing.T) {
	if d, err := NormalizeDomain(" Example.COM "); err != nil || d != "example.com" {
		t.Fatalf("got %q, %v", d, err)
	}
	for _, bad := range []string{"", "localhost", "a b.com", "-x.com", "http://x.com", "x..com"} {
		if _, err := NormalizeDomain(bad); err == nil {
			t.Errorf("%q should be rejected", bad)
		}
	}
}

func TestCloudflareOriginCertIsServedBySNI(t *testing.T) {
	e := newEngine(t)
	certPEM, keyPEM := selfSigned(t, "app.example.com")

	rec, err := e.Add(AddRequest{Domain: "app.example.com", Provider: ProviderCloudflareOrigin, CertPEM: certPEM, KeyPEM: keyPEM})
	if err != nil {
		t.Fatal(err)
	}
	if rec.Status != StatusActive {
		t.Fatalf("status = %s, want active", rec.Status)
	}
	c, err := e.GetCertificate(&tls.ClientHelloInfo{ServerName: "app.example.com"})
	if err != nil || c == nil {
		t.Fatalf("no certificate served for SNI: %v", err)
	}
	if _, err := e.GetCertificate(&tls.ClientHelloInfo{ServerName: "other.example.com"}); err == nil {
		t.Fatal("certificate served for an unrelated name")
	}
}

func TestCloudflareOriginRejectsWrongDomainAndMismatchedKey(t *testing.T) {
	e := newEngine(t)
	certPEM, keyPEM := selfSigned(t, "app.example.com")
	if _, err := e.Add(AddRequest{Domain: "other.example.com", Provider: ProviderCloudflareOrigin, CertPEM: certPEM, KeyPEM: keyPEM}); err == nil {
		t.Fatal("certificate for another domain was accepted")
	}
	otherCert, _ := selfSigned(t, "app.example.com")
	if _, err := e.Add(AddRequest{Domain: "app.example.com", Provider: ProviderCloudflareOrigin, CertPEM: otherCert, KeyPEM: keyPEM}); err == nil {
		t.Fatal("mismatched key was accepted")
	}
}

func TestZeroSSLNeedsEAB(t *testing.T) {
	e := newEngine(t)
	if _, err := e.Add(AddRequest{Domain: "app.example.com", Provider: ProviderZeroSSL}); err == nil {
		t.Fatal("ZeroSSL without EAB keys was accepted")
	}
}

func TestUnknownProviderRejected(t *testing.T) {
	e := newEngine(t)
	if _, err := e.Add(AddRequest{Domain: "app.example.com", Provider: "nope"}); err == nil {
		t.Fatal("unknown provider was accepted")
	}
}

func TestHTTP01ChallengeIsServed(t *testing.T) {
	e := newEngine(t)
	if err := e.Present("app.example.com", "tok123", "tok123.keyauth"); err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", challengePathPrefix+"tok123", nil)
	if !e.ServeChallenge(rec, req) {
		t.Fatal("challenge path not handled")
	}
	if rec.Body.String() != "tok123.keyauth" {
		t.Fatalf("body = %q", rec.Body.String())
	}
	_ = e.CleanUp("app.example.com", "tok123", "tok123.keyauth")
	rec = httptest.NewRecorder()
	e.ServeChallenge(rec, httptest.NewRequest("GET", challengePathPrefix+"tok123", nil))
	if rec.Code != 404 {
		t.Fatalf("cleaned token still served: %d", rec.Code)
	}
}

func TestRemoveStopsServing(t *testing.T) {
	e := newEngine(t)
	certPEM, keyPEM := selfSigned(t, "app.example.com")
	if _, err := e.Add(AddRequest{Domain: "app.example.com", Provider: ProviderCloudflareOrigin, CertPEM: certPEM, KeyPEM: keyPEM}); err != nil {
		t.Fatal(err)
	}
	if err := e.Remove("app.example.com"); err != nil {
		t.Fatal(err)
	}
	if _, err := e.GetCertificate(&tls.ClientHelloInfo{ServerName: "app.example.com"}); err == nil {
		t.Fatal("removed domain still served")
	}
}
