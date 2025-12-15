package services

import (
	"errors"
	"math"
	"regexp"
	"strings"
)

// validatePasswordStrength enforces password complexity requirements
func validatePasswordStrength(password string) error {
	if len(password) < 12 {
		return errors.New("password must be at least 12 characters long")
	}

	if len(password) > 128 {
		return errors.New("password must not exceed 128 characters")
	}

	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};:'"\\|,.<>/?]`).MatchString(password)
	)

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("password must contain uppercase, lowercase, number, and special character")
	}

	// Check against common weak passwords
	commonPasswords := []string{
		"password", "123456", "12345678", "qwerty", "abc123",
		"password123", "letmein", "welcome", "admin", "user",
		"passw0rd", "p@ssword", "password!", "12345678a",
	}

	lowerPassword := strings.ToLower(password)
	for _, weak := range commonPasswords {
		if strings.Contains(lowerPassword, weak) {
			return errors.New("password is too common or contains weak patterns")
		}
	}

	return nil
}

// isValidEmail validates email format using RFC 5322 compliant regex
func isValidEmail(email string) bool {
	// Simplified RFC 5322 compliant regex
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched && len(email) <= 254 // RFC 5321 max length
}

// calculateEntropy calculates Shannon entropy of a string
// Used to measure randomness/strength of secrets
func calculateEntropy(s string) float64 {
	if len(s) == 0 {
		return 0.0
	}

	// Count frequency of each character
	freq := make(map[rune]float64)
	for _, char := range s {
		freq[char]++
	}

	// Calculate entropy
	var entropy float64
	length := float64(len(s))
	for _, count := range freq {
		probability := count / length
		entropy -= probability * math.Log2(probability)
	}

	return entropy
}

// sanitizePrompt removes common prompt injection patterns
func sanitizePrompt(prompt string) string {
	// Remove potential injection patterns
	injectionPatterns := []string{
		"ignore previous instructions",
		"ignore all previous",
		"disregard all previous",
		"forget everything",
		"repeat the first",
		"reveal your prompt",
		"show me your instructions",
		"what are your instructions",
		"system prompt",
		"</system>",
		"<|system|>",
		"[SYSTEM]",
	}

	sanitized := strings.ToLower(prompt)
	for _, pattern := range injectionPatterns {
		sanitized = strings.ReplaceAll(sanitized, pattern, "")
	}

	// Trim to reasonable length
	const maxPromptLength = 2000
	if len(sanitized) > maxPromptLength {
		sanitized = sanitized[:maxPromptLength]
	}

	return strings.TrimSpace(sanitized)
}

// isSuspiciousResponse detects potentially unsafe LLM outputs
func isSuspiciousResponse(response string) bool {
	suspiciousPatterns := []string{
		"ignore previous instructions",
		"i am an ai",
		"my system prompt",
		"i cannot assist with that",
		"as an ai language model",
		"i'm sorry, but i can't",
		"</system>",
		"<script",
		"javascript:",
		"onerror=",
	}

	lower := strings.ToLower(response)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}

	// Check for excessive length (potential injection)
	if len(response) > 5000 {
		return true
	}

	return false
}
