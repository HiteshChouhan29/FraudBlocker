package utils

import (
	"regexp"
	"strings"
)

// NormalizePhoneNumber normalizes Indian phone numbers by removing non-digit characters and country code
func NormalizePhoneNumber(phone string) string {
	// Remove all non-digit characters
	re := regexp.MustCompile(`\D`)
	digits := re.ReplaceAllString(phone, "")

	// Remove leading '91' country code or leading '0'
	if strings.HasPrefix(digits, "91") && len(digits) > 10 {
		digits = digits[2:]
	} else if strings.HasPrefix(digits, "0") && len(digits) > 10 {
		digits = digits[1:]
	}

	return digits
}

// IsValidPhoneNumber validates Indian phone numbers (10 digits starting with 7,8, or 9)
func IsValidPhoneNumber(phone string) bool {
	normalized := NormalizePhoneNumber(phone)
	if len(normalized) != 10 {
		return false
	}
	firstDigit := normalized[0]
	return firstDigit == '7' || firstDigit == '8' || firstDigit == '9'
}

// FormatPhoneNumber formats a 10-digit Indian phone number as +91 XXXXX XXXXX
func FormatPhoneNumber(phone string) string {
	normalized := NormalizePhoneNumber(phone)
	if len(normalized) != 10 {
		return phone // Return original if not a standard 10-digit number
	}

	return "+91 " + normalized[:5] + " " + normalized[5:]
}
