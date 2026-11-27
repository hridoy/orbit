package s3

import (
	"encoding/json"
	"net/http"
	"strings"
)

// IAMPolicyDocument represents an AWS-compatible S3 IAM JSON policy.
type IAMPolicyDocument struct {
	Version   string               `json:"Version,omitempty"`
	Statement []IAMPolicyStatement `json:"Statement"`
}

type IAMPolicyStatement struct {
	Effect   string   `json:"Effect"`   // "Allow" or "Deny"
	Action   []string `json:"Action"`   // e.g. ["s3:GetObject", "s3:PutObject", "s3:*"]
	Resource []string `json:"Resource"` // e.g. ["arn:aws:s3:::assets/prod/*"]
}

func methodToS3Action(method, objectKey string) string {
	if objectKey == "" && method == http.MethodGet {
		return "s3:ListBucket"
	}
	switch method {
	case http.MethodGet, http.MethodHead:
		return "s3:GetObject"
	case http.MethodPut, http.MethodPost:
		return "s3:PutObject"
	case http.MethodDelete:
		return "s3:DeleteObject"
	default:
		return "s3:Unknown"
	}
}

func wildcardMatch(pattern, value string) bool {
	if pattern == "*" || pattern == value {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(value, strings.TrimSuffix(pattern, "*"))
	}
	return false
}

// evaluateIAMPolicyJSON evaluates an AWS-style JSON IAM policy document.
// Explicit Deny always overrides Allow. Default is Deny if no statement matches.
func evaluateIAMPolicyJSON(policyJSON, method, bucketName, objectKey string) (bool, error) {
	var doc IAMPolicyDocument
	if err := json.Unmarshal([]byte(policyJSON), &doc); err != nil {
		return false, err
	}

	action := methodToS3Action(method, objectKey)
	resource := "arn:aws:s3:::" + bucketName
	if objectKey != "" {
		resource += "/" + objectKey
	}

	allowed := false
	for _, stmt := range doc.Statement {
		actionMatch := false
		for _, a := range stmt.Action {
			if a == "*" || a == "s3:*" || wildcardMatch(strings.ToLower(a), strings.ToLower(action)) {
				actionMatch = true
				break
			}
		}
		if !actionMatch {
			continue
		}

		resMatch := false
		for _, r := range stmt.Resource {
			if r == "*" || wildcardMatch(r, resource) {
				resMatch = true
				break
			}
		}
		if !resMatch {
			continue
		}

		if strings.EqualFold(stmt.Effect, "Deny") {
			return false, nil
		}
		if strings.EqualFold(stmt.Effect, "Allow") {
			allowed = true
		}
	}
	return allowed, nil
}
