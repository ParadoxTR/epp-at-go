package validator

import (
	"fmt"
	"regexp"
)

var DomainNameRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)

func ValidateDomainName(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name cannot be empty")
	}

	if len(domain) > 253 {
		return fmt.Errorf("domain name too long: maximum 253 characters")
	}

	if !DomainNameRegex.MatchString(domain) {
		return fmt.Errorf("invalid domain name format: %s", domain)
	}

	return nil
}
