package helpers

import (
	"net/url"
	"strings"

	"github.com/gofrs/uuid"
)

// IsURL validates a full http/https URL
// source: https://stackoverflow.com/a/55551215
func IsURL(str string) bool {
	u, err := url.Parse(str)
	scheme := strings.ToLower(u.Scheme)
	return err == nil && (scheme == "http" || scheme == "https") && u.Host != ""
}

// IsHTTPURL returns true if it's a valid HTTP URL (no SSL)
func IsHTTPURL(str string) bool {
	u, err := url.Parse(str)
	scheme := strings.ToLower(u.Scheme)
	return err == nil && scheme == "http" && u.Host != ""
}

// RewriteToHTTPS replaces the protocol with https in a URL
func RewriteToHTTPS(str string) string {
	u, err := url.Parse(str)
	if err != nil {
		return str
	}
	u.Scheme = "https"
	return u.String()
}

// IsSameUUID returns true if both uuid are the same
func IsSameUUID(a, b uuid.UUID) bool {
	return a.String() == b.String()
}
