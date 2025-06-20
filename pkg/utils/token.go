package utils

import (
	"strings"
)

// ExtractTokenFromHeader extracts the JWT token from the Authorization header.
// It expects the header to be in the format: "Bearer <token>"
func ExtractTokenFromHeader(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}
