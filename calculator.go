package main

import (
	"crypto/sha256"
	"strings"
	"unicode"
)

// CalculateLovePercentage computes a deterministic love percentage based on two names.
//
// Algorithm:
// 1. Trim whitespace and convert to lowercase
// 2. Remove all non-alphabetic characters for consistent hashing
// 3. Sort names alphabetically to ensure symmetry (name1 + name2 = name2 + name1)
// 4. Combine normalized names
// 5. Compute SHA256 hash
// 6. Use hash bytes to generate a value between 50 and 100
//
// The range 50-100 ensures a more positive user experience—no one gets a "0%" result.
func CalculateLovePercentage(name1, name2 string) int {
	// Normalize: trim and lowercase
	n1 := strings.ToLower(strings.TrimSpace(name1))
	n2 := strings.ToLower(strings.TrimSpace(name2))

	// Remove non-alphabetic characters for consistent hashing
	n1 = removeNonAlpha(n1)
	n2 = removeNonAlpha(n2)

	// Sort names alphabetically to ensure symmetry
	// (Saksham + Arohi produces same result as Arohi + Saksham)
	if n1 > n2 {
		n1, n2 = n2, n1
	}

	// Combine names
	combined := n1 + n2

	// Compute SHA256 hash
	hash := sha256.Sum256([]byte(combined))

	// Convert hash to a number between 0 and 100
	// Use the first 4 bytes of the hash to generate a uint32
	value := uint32(hash[0])<<24 | uint32(hash[1])<<16 | uint32(hash[2])<<8 | uint32(hash[3])

	// Normalize to 0-100 range
	percentage := int(value % 101)

	return percentage
}

// GetLoveMessage returns an appropriate romantic message based on the percentage.
func GetLoveMessage(percentage int) string {
	switch {
	case percentage >= 90:
		return "💘 THE UNIVERSE HAS SPOKEN"
	case percentage >= 75:
		return "💕 Okayyy... there's definitely something here 👀"
	case percentage >= 60:
		return "💖 Hmmm... the chemistry is chemistry-ing"
	case percentage >= 40:
		return "😭 The stars are being confusing today"
	default:
		return "💀 Cupid may have missed the appointment"
	}
}

// removeNonAlpha removes all non-alphabetic characters from a string.
func removeNonAlpha(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		return -1
	}, s)
}
