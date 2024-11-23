package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// IsValidCPF validates a CPF number.
func IsValidCPF(cpf string) bool {
	// Remove non-numeric characters
	cpf = strings.TrimSpace(cpf)
	re := regexp.MustCompile(`\D`)
	cpf = re.ReplaceAllString(cpf, "")

	// CPF must have 11 digits
	if len(cpf) != 11 {
		return false
	}

	// Reject sequences like "111.111.111-11"
	for i := 0; i < 10; i++ {
		if cpf == strings.Repeat(strconv.Itoa(i), 11) {
			return false
		}
	}

	// Validate first check digit
	if !validateDigit(cpf, 9) {
		return false
	}

	// Validate second check digit
	return validateDigit(cpf, 10)
}

// validateDigit calculates and validates a CPF check digit.
func validateDigit(cpf string, position int) bool {
	sum := 0
	for i := 0; i < position; i++ {
		num, _ := strconv.Atoi(string(cpf[i]))
		sum += num * (position + 1 - i)
	}

	remainder := sum % 11
	expectedDigit := 0
	if remainder >= 2 {
		expectedDigit = 11 - remainder
	}

	actualDigit, _ := strconv.Atoi(string(cpf[position]))
	return actualDigit == expectedDigit
}
