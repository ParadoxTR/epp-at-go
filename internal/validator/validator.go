package validator

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
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

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email address cannot be empty")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email address: %s", email)
	}

	return nil
}

func ValidateContactID(contactID string) error {
	if contactID == "" {
		return fmt.Errorf("contact ID cannot be empty")
	}

	if contactID == "AUTO" {
		return nil
	}

	if len(contactID) < 3 || len(contactID) > 16 {
		return fmt.Errorf("contact ID must be between 3 and 16 characters")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`).MatchString(contactID) {
		return fmt.Errorf("contact ID can only contain alphanumeric characters, hyphens, and underscores")
	}

	return nil
}

func ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return nil // Phone number is optional in many cases
	}

	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")

	if !regexp.MustCompile(`^\+[1-9]\d{1,14}$`).MatchString(cleaned) {
		return fmt.Errorf("invalid phone number format: must be in international format (+country.number)")
	}

	return nil
}

func ValidateCountryCode(cc string) error {
	if cc == "" {
		return fmt.Errorf("country code cannot be empty")
	}

	if len(cc) != 2 {
		return fmt.Errorf("country code must be exactly 2 characters")
	}

	if !regexp.MustCompile(`^[A-Z]{2}$`).MatchString(cc) {
		return fmt.Errorf("country code must be two uppercase letters")
	}

	return nil
}

func ValidatePostalCode(pc string) error {
	if pc == "" {
		return fmt.Errorf("postal code cannot be empty")
	}

	if len(pc) > 16 {
		return fmt.Errorf("postal code too long: maximum 16 characters")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9\s\-]+$`).MatchString(pc) {
		return fmt.Errorf("postal code contains invalid characters")
	}

	return nil
}

func ValidateAuthInfo(authInfo string) error {
	if authInfo == "" {
		return fmt.Errorf("auth info cannot be empty")
	}

	if len(authInfo) < 6 {
		return fmt.Errorf("auth info must be at least 6 characters long")
	}

	if len(authInfo) > 64 {
		return fmt.Errorf("auth info too long: maximum 64 characters")
	}

	return nil
}
