// Package tlsengine issues and renews HTTPS certificates for public domains. Each domain
// picks a provider: Let's Encrypt or ZeroSSL (ACME, HTTP-01 on port 80), or a Cloudflare
// Origin Certificate pasted in by the operator. Certificates are served by SNI on port 443.
package tlsengine

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"

	"orbit/internal/db"
	"orbit/internal/security"
)

const (
	ProviderLetsEncrypt      = "letsencrypt"
	ProviderZeroSSL          = "zerossl"
	ProviderCloudflareOrigin = "cloudflare_origin"

	StatusPending = "pending"
	StatusActive  = "active"
	StatusFailed  = "failed"

	challengePathPrefix = "/.well-known/acme-challenge/"
	renewBefore         = 30 * 24 * time.Hour
	renewEvery          = 6 * time.Hour
	zerosslDirectory    = "https://acme.zerossl.com/v2/DV90"
)

var domainRe = regexp.MustCompile(`^([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,63}$`)

// List returns every configured domain, newest changes first.
func (e *Engine) List() ([]*db.TLSDomain, error) {
	return e.store.List()
}

// ValidProvider reports whether p is a supported provider name.
func ValidProvider(p string) bool {
	return p == ProviderLetsEncrypt || p == ProviderZeroSSL || p == ProviderCloudflareOrigin
}

// NormalizeDomain lowercases a domain and rejects anything that is not a plain hostname.
func NormalizeDomain(d string) (string, error) {
	d = strings.ToLower(strings.TrimSpace(d))
	if !domainRe.MatchString(d) {
		return "", errors.New("enter a plain domain name such as example.com")
	}
	return d, nil
}

// AddRequest is what the operator submits for one domain.
type AddRequest struct {
	Domain   string `json:"domain"`
	Provider string `json:"provider"`
	EABKID   string `json:"eab_kid,omitempty"`      // ZeroSSL only
	EABHMAC  string `json:"eab_hmac,omitempty"`     // ZeroSSL only
	CertPEM  string `json:"cert_pem,omitempty"`     // Cloudflare origin only
	KeyPEM   string `json:"key_pem,omitempty"`      // Cloudflare origin only
	CFToken  string `json:"cf_api_token,omitempty"` // Cloudflare origin only: Orbit issues and renews it
}

// Engine holds the live certificates and runs issuance and renewal.
type Engine struct {
	store *db.TLSDomainStore
	vault *security.Vault

	mu    sync.RWMutex
	certs map[string]*tls.Certificate // by domain

	challenges sync.Map // HTTP-01 token -> key authorization

	// directory overrides the ACME directory for Let's Encrypt (tests only).
	directory string
}

func New(store *db.TLSDomainStore, vault *security.Vault) *Engine {
	return &Engine{store: store, vault: vault, certs: map[string]*tls.Certificate{}, directory: lego.LEDirectoryProduction}
}

// Load reads every active certificate from the database into memory.
func (e *Engine) Load() error {
	list, err := e.store.List()
	if err != nil {
		return err
	}
	for _, d := range list {
		if d.Status != StatusActive {
			continue
		}
		cert, err := e.parseStored(d)
		if err != nil {
			log.Printf("⚠️ tls: %s: %v", d.Domain, err)
			continue
		}
		e.setCert(d.Domain, cert)
	}
	return nil
}

// GetCertificate picks the certificate for the SNI name in the TLS handshake.
func (e *Engine) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if c, ok := e.certs[strings.ToLower(hello.ServerName)]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("no certificate for %q", hello.ServerName)
}

// HasCertificate reports whether a certificate is loaded for the domain.
func (e *Engine) HasCertificate(domain string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	_, ok := e.certs[strings.ToLower(domain)]
	return ok
}

// ServeChallenge answers ACME HTTP-01 requests. It returns true when the request was handled.
func (e *Engine) ServeChallenge(w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.URL.Path, challengePathPrefix) {
		return false
	}
	token := strings.TrimPrefix(r.URL.Path, challengePathPrefix)
	if v, ok := e.challenges.Load(token); ok {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(v.(string)))
		return true
	}
	http.NotFound(w, r)
	return true
}

// Add validates a domain and starts its certificate. ACME providers issue in the background
// and report progress through the status field. A Cloudflare origin certificate is active at once.
func (e *Engine) Add(req AddRequest) (*db.TLSDomain, error) {
	domain, err := NormalizeDomain(req.Domain)
	if err != nil {
		return nil, err
	}
	if !ValidProvider(req.Provider) {
		return nil, errors.New("unknown provider")
	}

	rec := &db.TLSDomain{Domain: domain, Provider: req.Provider, Status: StatusPending}

	switch req.Provider {
	case ProviderCloudflareOrigin:
		if req.CertPEM == "" && req.KeyPEM == "" {
			if req.CFToken == "" {
				return nil, errors.New("enter a Cloudflare API token, or paste the origin certificate and key")
			}
			return e.addCloudflareAuto(domain, req.CFToken)
		}
		cert, err := tls.X509KeyPair([]byte(req.CertPEM), []byte(req.KeyPEM))
		if err != nil {
			return nil, fmt.Errorf("certificate and key do not match: %w", err)
		}
		leaf, err := leafOf(&cert)
		if err != nil {
			return nil, err
		}
		if err := leaf.VerifyHostname(domain); err != nil {
			return nil, fmt.Errorf("certificate is not for %s: %w", domain, err)
		}
		keyEnc, err := e.seal([]byte(req.KeyPEM))
		if err != nil {
			return nil, err
		}
		exp := leaf.NotAfter
		rec.Status = StatusActive
		rec.CertPEM = req.CertPEM
		rec.KeyEnc = keyEnc
		rec.ExpiresAt = &exp
		if err := e.store.Upsert(rec); err != nil {
			return nil, err
		}
		e.setCert(domain, &cert)
		return rec, nil

	default: // ACME providers
		if req.Provider == ProviderZeroSSL && (req.EABKID == "" || req.EABHMAC == "") {
			return nil, errors.New("ZeroSSL needs the EAB key ID and HMAC key from your ZeroSSL account")
		}
		if req.Provider == ProviderZeroSSL {
			hmacEnc, err := e.seal([]byte(req.EABHMAC))
			if err != nil {
				return nil, err
			}
			rec.EABKID = req.EABKID
			rec.EABHMACEnc = hmacEnc
		}
		if err := e.store.Upsert(rec); err != nil {
			return nil, err
		}
		go e.issue(domain, req.Provider, req.EABKID, req.EABHMAC)
		return rec, nil
	}
}

// addCloudflareAuto issues an Origin CA certificate through the Cloudflare API. The token is
// stored encrypted so renewals can run without the operator.
func (e *Engine) addCloudflareAuto(domain, token string) (*db.TLSDomain, error) {
	certPEM, keyPEM, exp, err := issueCloudflareOrigin(context.Background(), token, domain)
	if err != nil {
		return nil, err
	}
	tc, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}
	keyEnc, err := e.seal(keyPEM)
	if err != nil {
		return nil, err
	}
	tokEnc, err := e.seal([]byte(token))
	if err != nil {
		return nil, err
	}
	rec := &db.TLSDomain{
		Domain: domain, Provider: ProviderCloudflareOrigin, Status: StatusActive,
		CertPEM: string(certPEM), KeyEnc: keyEnc, EABHMACEnc: tokEnc, ExpiresAt: &exp,
	}
	if err := e.store.Upsert(rec); err != nil {
		return nil, err
	}
	e.setCert(domain, &tc)
	return rec, nil
}

// renewCloudflareAuto reissues a Cloudflare origin certificate that was created with a token.
func (e *Engine) renewCloudflareAuto(d *db.TLSDomain) {
	token, err := e.open(d.EABHMACEnc)
	if err != nil {
		return
	}
	rec, err := e.addCloudflareAuto(d.Domain, string(token))
	if err != nil {
		if cur, gerr := e.store.Get(d.Domain); gerr == nil && cur != nil {
			cur.LastError = err.Error()
			_ = e.store.Upsert(cur)
		}
		log.Printf("⚠️ tls: renewing %s from Cloudflare failed: %v", d.Domain, err)
		return
	}
	_ = rec
}

// Remove deletes a domain and its certificate and stops serving it.
func (e *Engine) Remove(domain string) error {
	d, err := NormalizeDomain(domain)
	if err != nil {
		return err
	}
	e.mu.Lock()
	delete(e.certs, d)
	e.mu.Unlock()
	return e.store.Delete(d)
}

// issue obtains a certificate from an ACME provider and records the result.
func (e *Engine) issue(domain, provider, kid, hmacKey string) {
	cert, certPEM, keyPEM, err := e.obtain(domain, provider, kid, hmacKey)
	rec, gerr := e.store.Get(domain)
	if gerr != nil || rec == nil {
		return
	}
	if err != nil {
		rec.Status = StatusFailed
		rec.LastError = err.Error()
		_ = e.store.Upsert(rec)
		log.Printf("⚠️ tls: issuing %s via %s failed: %v", domain, provider, err)
		return
	}
	keyEnc, serr := e.seal(keyPEM)
	if serr != nil {
		rec.Status = StatusFailed
		rec.LastError = serr.Error()
		_ = e.store.Upsert(rec)
		return
	}
	exp := cert.NotAfter
	rec.Status = StatusActive
	rec.LastError = ""
	rec.CertPEM = string(certPEM)
	rec.KeyEnc = keyEnc
	rec.ExpiresAt = &exp
	if err := e.store.Upsert(rec); err != nil {
		log.Printf("⚠️ tls: saving %s: %v", domain, err)
		return
	}
	tc, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return
	}
	e.setCert(domain, &tc)
	log.Printf("🔒 tls: %s certificate active via %s until %s", domain, provider, exp.Format(time.DateOnly))
}

func (e *Engine) obtain(domain, provider, kid, hmacKey string) (*x509.Certificate, []byte, []byte, error) {
	account, err := newAccount()
	if err != nil {
		return nil, nil, nil, err
	}
	cfg := lego.NewConfig(account)
	cfg.CADirURL = e.directory
	if provider == ProviderZeroSSL {
		cfg.CADirURL = zerosslDirectory
	}
	cfg.Certificate.KeyType = certcrypto.EC256

	client, err := lego.NewClient(cfg)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := client.Challenge.SetHTTP01Provider(e); err != nil {
		return nil, nil, nil, err
	}

	if provider == ProviderZeroSSL {
		_, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
			TermsOfServiceAgreed: true,
			Kid:                  kid,
			HmacEncoded:          hmacKey,
		})
	} else {
		_, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	}
	if err != nil {
		return nil, nil, nil, fmt.Errorf("account registration: %w", err)
	}

	res, err := client.Certificate.Obtain(certificate.ObtainRequest{Domains: []string{domain}, Bundle: true})
	if err != nil {
		return nil, nil, nil, err
	}
	leaf, err := parseLeafPEM(res.Certificate)
	if err != nil {
		return nil, nil, nil, err
	}
	return leaf, res.Certificate, res.PrivateKey, nil
}

// RunRenewals checks every ACME certificate about every six hours and reissues the ones close to expiry.
func (e *Engine) RunRenewals(ctx context.Context) {
	ticker := time.NewTicker(renewEvery)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			list, err := e.store.List()
			if err != nil {
				continue
			}
			for _, d := range list {
				if d.Status != StatusActive || d.ExpiresAt == nil {
					continue
				}
				if time.Until(*d.ExpiresAt) > renewBefore {
					continue
				}
				if d.Provider == ProviderCloudflareOrigin {
					if d.EABHMACEnc != "" { // created with a token, so it can renew itself
						go e.renewCloudflareAuto(d)
					}
					continue
				}
				kid, hmacKey := "", ""
				if d.Provider == ProviderZeroSSL {
					kid = d.EABKID
					raw, err := e.open(d.EABHMACEnc)
					if err != nil {
						continue
					}
					hmacKey = string(raw)
				}
				go e.issue(d.Domain, d.Provider, kid, hmacKey)
			}
		}
	}
}

// Present and CleanUp implement lego's HTTP-01 provider. The token is served by ServeChallenge.
func (e *Engine) Present(domain, token, keyAuth string) error {
	e.challenges.Store(token, keyAuth)
	return nil
}

func (e *Engine) CleanUp(domain, token, keyAuth string) error {
	e.challenges.Delete(token)
	return nil
}

func (e *Engine) setCert(domain string, c *tls.Certificate) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.certs[domain] = c
}

func (e *Engine) parseStored(d *db.TLSDomain) (*tls.Certificate, error) {
	keyPEM, err := e.open(d.KeyEnc)
	if err != nil {
		return nil, err
	}
	c, err := tls.X509KeyPair([]byte(d.CertPEM), keyPEM)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// seal encrypts a secret for storage. The format is ciphertext:nonce, both base64.
func (e *Engine) seal(plain []byte) (string, error) {
	ct, nonce, err := e.vault.Encrypt(plain)
	if err != nil {
		return "", err
	}
	return ct + ":" + nonce, nil
}

func (e *Engine) open(enc string) ([]byte, error) {
	parts := strings.SplitN(enc, ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("stored secret is malformed")
	}
	return e.vault.Decrypt(parts[0], parts[1])
}

// accountUser is an ACME account with a fresh key. A new key per issuance keeps no account
// secrets on disk.
type accountUser struct {
	key crypto.PrivateKey
}

func newAccount() (*accountUser, error) {
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &accountUser{key: k}, nil
}

func (a *accountUser) GetEmail() string                        { return "" }
func (a *accountUser) GetRegistration() *registration.Resource { return nil }
func (a *accountUser) GetPrivateKey() crypto.PrivateKey        { return a.key }

func leafOf(c *tls.Certificate) (*x509.Certificate, error) {
	if c.Leaf != nil {
		return c.Leaf, nil
	}
	return x509.ParseCertificate(c.Certificate[0])
}

func parseLeafPEM(certPEM []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, errors.New("certificate PEM not found")
	}
	return x509.ParseCertificate(block.Bytes)
}
