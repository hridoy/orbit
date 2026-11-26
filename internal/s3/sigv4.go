package s3

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	ErrMissingAuthHeader = errors.New("missing authorization header or presigned query params")
	ErrInvalidSignature  = errors.New("the request signature we calculated does not match the signature you provided")
	ErrSignatureExpired  = errors.New("presigned URL has expired")
	ErrRequestTooSkewed  = errors.New("request timestamp is outside the allowed clock skew")
	ErrUnsignedWrite     = errors.New("UNSIGNED-PAYLOAD is not allowed for write requests")
	ErrPresignedWrite    = errors.New("presigned URLs are only accepted for GET and HEAD")
)

// maxClockSkew bounds how far a signed request's X-Amz-Date may drift from server time.
const maxClockSkew = 15 * time.Minute

var sha256HexRe = regexp.MustCompile(`^[0-9a-f]{64}$`)

type SigV4Auth struct {
	AccessKey string
	SecretKey string
	Region    string
}

func NewSigV4Auth(accessKey, secretKey, region string) *SigV4Auth {
	if region == "" {
		region = "us-east-1"
	}
	return &SigV4Auth{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Region:    region,
	}
}

// DeriveSigningKey computes kSigning = HMAC(HMAC(HMAC(HMAC("AWS4" + Secret, Date), Region), Service), "aws4_request")
func DeriveSigningKey(secretKey, dateStamp, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secretKey), []byte(dateStamp))
	kRegion := hmacSHA256(kDate, []byte(region))
	kService := hmacSHA256(kRegion, []byte(service))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))
	return kSigning
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// CanonicalQueryString sorts and escapes query parameters per AWS spec.
func CanonicalQueryString(query url.Values) string {
	if len(query) == 0 {
		return ""
	}

	keys := make([]string, 0, len(query))
	for k := range query {
		if k == "X-Amz-Signature" {
			continue // Exclude signature from canonical query
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		values := query[k]
		sort.Strings(values)
		for _, v := range values {
			parts = append(parts, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
		}
	}
	return strings.Join(parts, "&")
}

// checkClockSkew rejects requests whose X-Amz-Date is missing, malformed, or too far from now.
func checkClockSkew(amzDate string, now time.Time) error {
	t, err := time.Parse("20060102T150405Z", amzDate)
	if err != nil {
		return errors.New("missing or invalid X-Amz-Date")
	}
	drift := now.Sub(t)
	if drift < 0 {
		drift = -drift
	}
	if drift > maxClockSkew {
		return ErrRequestTooSkewed
	}
	return nil
}

// VerifyRequest authenticates an incoming S3 request (header or presigned query).
// On success it returns the payload SHA-256 (lowercase hex) that the request body
// must match, or "" when the body is not covered by the signature (presigned GET/HEAD
// and unsigned reads). Callers must verify the body against that hash for uploads.
func (v *SigV4Auth) VerifyRequest(r *http.Request) (string, error) {
	// 1. Presigned URL query params
	if r.URL.Query().Get("X-Amz-Signature") != "" {
		return "", v.verifyPresigned(r)
	}

	// 2. Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrMissingAuthHeader
	}
	if !strings.HasPrefix(authHeader, "AWS4-HMAC-SHA256 ") {
		return "", errors.New("unsupported authorization scheme")
	}

	authParams := strings.TrimPrefix(authHeader, "AWS4-HMAC-SHA256 ")
	parts := strings.Split(authParams, ", ")
	var credential, signedHeaders, signature string
	for _, p := range parts {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "Credential":
			credential = kv[1]
		case "SignedHeaders":
			signedHeaders = kv[1]
		case "Signature":
			signature = kv[1]
		}
	}
	if credential == "" || signedHeaders == "" || signature == "" {
		return "", errors.New("malformed authorization header")
	}

	// Credential format: <accessKey>/<dateStamp>/<region>/s3/aws4_request
	credParts := strings.Split(credential, "/")
	if len(credParts) != 5 {
		return "", errors.New("invalid credential scope in authorization header")
	}
	if credParts[0] != v.AccessKey {
		return "", errors.New("invalid access key ID")
	}
	dateStamp := credParts[1]
	reqRegion := credParts[2]

	amzDate := r.Header.Get("X-Amz-Date")
	if err := checkClockSkew(amzDate, time.Now().UTC()); err != nil {
		return "", err
	}

	// The payload hash must be declared and signed. Only a real hash is accepted for writes.
	payloadHash := r.Header.Get("X-Amz-Content-Sha256")
	switch {
	case payloadHash == "UNSIGNED-PAYLOAD":
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			return "", ErrUnsignedWrite
		}
	case sha256HexRe.MatchString(payloadHash):
	default:
		return "", errors.New("X-Amz-Content-Sha256 must be a SHA-256 hex digest")
	}

	var canonicalHeaders strings.Builder
	for _, h := range strings.Split(signedHeaders, ";") {
		hLower := strings.ToLower(strings.TrimSpace(h))
		val := strings.TrimSpace(r.Header.Get(hLower))
		canonicalHeaders.WriteString(fmt.Sprintf("%s:%s\n", hLower, val))
	}

	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		r.Method,
		r.URL.Path,
		CanonicalQueryString(r.URL.Query()),
		canonicalHeaders.String(),
		signedHeaders,
		payloadHash,
	)

	if err := v.checkSignature(canonicalRequest, amzDate, dateStamp, reqRegion, signature); err != nil {
		return "", err
	}

	if payloadHash == "UNSIGNED-PAYLOAD" {
		return "", nil
	}
	return payloadHash, nil
}

func (v *SigV4Auth) checkSignature(canonicalRequest, amzDate, dateStamp, region, signature string) error {
	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request", dateStamp, region)
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		amzDate,
		credentialScope,
		sha256Hex([]byte(canonicalRequest)),
	)
	signingKey := DeriveSigningKey(v.SecretKey, dateStamp, region, "s3")
	expected := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))
	if subtle.ConstantTimeCompare([]byte(signature), []byte(expected)) != 1 {
		return ErrInvalidSignature
	}
	return nil
}

func (v *SigV4Auth) verifyPresigned(r *http.Request) error {
	// Presigned URLs carry no body hash, so they are only valid for reads.
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return ErrPresignedWrite
	}

	q := r.URL.Query()
	signature := q.Get("X-Amz-Signature")
	credential := q.Get("X-Amz-Credential")
	amzDate := q.Get("X-Amz-Date")
	expiresStr := q.Get("X-Amz-Expires")
	signedHeaders := q.Get("X-Amz-SignedHeaders")

	reqTime, err := time.Parse("20060102T150405Z", amzDate)
	if err != nil {
		return errors.New("invalid X-Amz-Date format")
	}

	var expiresSec int64 = 900
	if expiresStr != "" {
		_, _ = fmt.Sscanf(expiresStr, "%d", &expiresSec)
	}
	if expiresSec <= 0 || expiresSec > 7*24*3600 {
		return errors.New("X-Amz-Expires must be between 1 and 604800 seconds")
	}

	now := time.Now().UTC()
	if now.After(reqTime.Add(time.Duration(expiresSec) * time.Second)) {
		return ErrSignatureExpired
	}
	if reqTime.After(now.Add(maxClockSkew)) {
		return ErrRequestTooSkewed
	}

	credParts := strings.Split(credential, "/")
	if len(credParts) != 5 || credParts[0] != v.AccessKey {
		return errors.New("invalid presigned credential")
	}
	dateStamp := credParts[1]
	reqRegion := credParts[2]

	var canonicalHeaders strings.Builder
	for _, h := range strings.Split(signedHeaders, ";") {
		hLower := strings.ToLower(strings.TrimSpace(h))
		val := strings.TrimSpace(r.Header.Get(hLower))
		canonicalHeaders.WriteString(fmt.Sprintf("%s:%s\n", hLower, val))
	}

	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\nUNSIGNED-PAYLOAD",
		r.Method,
		r.URL.Path,
		CanonicalQueryString(q),
		canonicalHeaders.String(),
		signedHeaders,
	)

	return v.checkSignature(canonicalRequest, amzDate, dateStamp, reqRegion, signature)
}

// GeneratePresignedURL creates an AWS SigV4 presigned URL for GET operations.
func (v *SigV4Auth) GeneratePresignedURL(method, host, bucket, key string, ttl time.Duration) string {
	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")

	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request", dateStamp, v.Region)
	credential := fmt.Sprintf("%s/%s", v.AccessKey, credentialScope)

	path := fmt.Sprintf("/s3/%s/%s", bucket, strings.TrimPrefix(key, "/"))

	q := url.Values{}
	q.Set("X-Amz-Algorithm", "AWS4-HMAC-SHA256")
	q.Set("X-Amz-Credential", credential)
	q.Set("X-Amz-Date", amzDate)
	q.Set("X-Amz-Expires", fmt.Sprintf("%d", int(ttl.Seconds())))
	q.Set("X-Amz-SignedHeaders", "host")

	canonicalHeaders := fmt.Sprintf("host:%s\n", host)
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\nhost\nUNSIGNED-PAYLOAD",
		method,
		path,
		CanonicalQueryString(q),
		canonicalHeaders,
	)

	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		amzDate,
		credentialScope,
		sha256Hex([]byte(canonicalRequest)),
	)

	signingKey := DeriveSigningKey(v.SecretKey, dateStamp, v.Region, "s3")
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	q.Set("X-Amz-Signature", signature)
	return fmt.Sprintf("http://%s%s?%s", host, path, q.Encode())
}
