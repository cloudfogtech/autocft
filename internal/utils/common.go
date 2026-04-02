package utils

import "strings"

const PerPage = 1000.0

func NormalizeHostname(hostname string) string {
	hostname = strings.TrimSpace(hostname)
	hostname = strings.TrimSuffix(hostname, ".")
	return strings.ToLower(hostname)
}
