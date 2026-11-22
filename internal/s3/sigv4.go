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
	"sort"
	"strings"
	"time"
)

var (
	ErrMissingAuthHeader = errors.New("missing authorization header or presigned query params")
	ErrInvalidSignature  = errors.New("the request signature we calculated does not match the signature you provided")
	ErrSignatureExpired  = errors.New("presigned URL has expired")
)

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

// VerifyRequest authenticates an incoming S3 request (either via Header or Presigned Query).
func (v *SigV4Auth) VerifyRequest(r *http.Request, payloadHash string) error {
	// 1. Check Presigned URL query params first
	q := r.URL.Query()
	if q.Get("X-Amz-Signature") != "" {
		return v.verifyPresigned(r)
	}

	// 2. Check Authorization Header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ErrMissingAuthHeader
	}

	if !strings.HasPrefix(authHeader, "AWS4-HMAC-SHA256 ") {
		return errors.New("unsupported authorization scheme")
	}

	// Parse header components
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
		return errors.New("malformed authorization header")
	}

	// Credential format: <accessKey>/<dateStamp>/<region>/s3/aws4_request
	credParts := strings.Split(credential, "/")
	if len(credParts) != 5 {
		return errors.New("invalid credential scope in authorization header")
	}

	reqAccessKey := credParts[0]
	dateStamp := credParts[1]
	reqRegion := credParts[2]

	if reqAccessKey != v.AccessKey {
		return errors.New("invalid access key ID")
	}

	// Determine request date from X-Amz-Date or Date
	amzDate := r.Header.Get("X-Amz-Date")
	if amzDate == "" {
		amzDate = r.Header.Get("Date")
	}

	// Canonical headers and signed headers list
	shList := strings.Split(signedHeaders, ";")
	var canonicalHeaders strings.Builder
	for _, h := range shList {
		hLower := strings.ToLower(strings.TrimSpace(h))
		val := strings.TrimSpace(r.Header.Get(hLower))
		canonicalHeaders.WriteString(fmt.Sprintf("%s:%s\n", hLower, val))
	}

	if payloadHash == "" {
		payloadHash = r.Header.Get("X-Amz-Content-Sha256")
		if payloadHash == "" {
			payloadHash = "UNSIGNED-PAYLOAD"
		}
	}

	// Build Canonical Request
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		r.Method,
		r.URL.Path,
		CanonicalQueryString(r.URL.Query()),
		canonicalHeaders.String(),
		signedHeaders,
		payloadHash,
	)

	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request", dateStamp, reqRegion)
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		amzDate,
		credentialScope,
		sha256Hex([]byte(canonicalRequest)),
	)

	signingKey := DeriveSigningKey(v.SecretKey, dateStamp, reqRegion, "s3")
	expectedSignature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	if subtle.ConstantTimeCompare([]byte(signature), []byte(expectedSignature)) != 1 {
		return ErrInvalidSignature
	}

	return nil
}

func (v *SigV4Auth) verifyPresigned(r *http.Request) error {
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

	if time.Now().UTC().After(reqTime.Add(time.Duration(expiresSec) * time.Second)) {
		return ErrSignatureExpired
	}

	credParts := strings.Split(credential, "/")
	if len(credParts) != 5 || credParts[0] != v.AccessKey {
		return errors.New("invalid presigned credential")
	}
	dateStamp := credParts[1]
	reqRegion := credParts[2]

	shList := strings.Split(signedHeaders, ";")
	var canonicalHeaders strings.Builder
	for _, h := range shList {
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

	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request", dateStamp, reqRegion)
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		amzDate,
		credentialScope,
		sha256Hex([]byte(canonicalRequest)),
	)

	signingKey := DeriveSigningKey(v.SecretKey, dateStamp, reqRegion, "s3")
	expectedSignature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	if subtle.ConstantTimeCompare([]byte(signature), []byte(expectedSignature)) != 1 {
		return ErrInvalidSignature
	}

	return nil
}

// GeneratePresignedURL creates an AWS SigV4 presigned URL for GET/PUT operations.
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
