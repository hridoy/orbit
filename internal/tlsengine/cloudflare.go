package tlsengine

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var cloudflareOriginAPI = "https://api.cloudflare.com/client/v4/certificates"

// originValidityDays is the lifetime asked of Cloudflare. Renewal runs well before expiry.
const originValidityDays = 5475

// issueCloudflareOrigin asks Cloudflare to sign an Origin CA certificate for domain. The key
// is made here and never leaves this server. The API token needs Zone > SSL and Certificates > Edit.
func issueCloudflareOrigin(ctx context.Context, apiToken, domain string) (certPEM, keyPEM []byte, expires time.Time, err error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, time.Time{}, err
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: domain},
		DNSNames: []string{domain},
	}, key)
	if err != nil {
		return nil, nil, time.Time{}, err
	}
	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})

	body, _ := json.Marshal(map[string]any{
		"hostnames":          []string{domain},
		"requested_validity": originValidityDays,
		"request_type":       "origin-ecc",
		"csr":                string(csrPEM),
	})
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cloudflareOriginAPI, bytes.NewReader(body))
	if err != nil {
		return nil, nil, time.Time{}, err
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, time.Time{}, err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))

	var out struct {
		Success bool `json:"success"`
		Errors  []struct {
			Message string `json:"message"`
		} `json:"errors"`
		Result struct {
			Certificate string `json:"certificate"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, nil, time.Time{}, fmt.Errorf("cloudflare returned an unexpected response (HTTP %d)", res.StatusCode)
	}
	if !out.Success || out.Result.Certificate == "" {
		msg := fmt.Sprintf("HTTP %d", res.StatusCode)
		if len(out.Errors) > 0 {
			msg = out.Errors[0].Message
		}
		return nil, nil, time.Time{}, errors.New("cloudflare refused the certificate: " + msg)
	}

	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, time.Time{}, err
	}
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	leaf, err := parseLeafPEM([]byte(out.Result.Certificate))
	if err != nil {
		return nil, nil, time.Time{}, err
	}
	return []byte(out.Result.Certificate), keyPEM, leaf.NotAfter, nil
}
