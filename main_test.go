package main

import (
	"os"
	"testing"
)

func TestCalculateLovePercentage(t *testing.T) {
	// Test symmetry
	p1 := CalculateLovePercentage("Alice", "Bob")
	p2 := CalculateLovePercentage("Bob", "Alice")
	if p1 != p2 {
		t.Errorf("Expected symmetric results, got %d and %d", p1, p2)
	}

	// Test bounds
	if p1 < 0 || p1 > 100 {
		t.Errorf("Percentage %d is out of bounds (0-100)", p1)
	}

	// Test special case: Saksham + Kirtika should always be 100%
	testCases := []struct {
		name1 string
		name2 string
	}{
		{"Saksham", "Kirtika"},
		{"Kirtika", "Saksham"},
		{"SAKSHAM", "KIRTIKA"},
		{"kirtika", "saksham"},
		{"SaKsHaM", "KiRtIkA"},
	}

	for _, tc := range testCases {
		result := CalculateLovePercentage(tc.name1, tc.name2)
		if result != 100 {
			t.Errorf("Expected 100%% for %s + %s, got %d%%", tc.name1, tc.name2, result)
		}
	}
}

func TestLocalStore(t *testing.T) {
	testFile := "test_submissions.json"
	defer os.Remove(testFile)

	s := newFileStore(testFile)
	err := s.Save("Romeo", "Juliet", 99)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	subs, err := s.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(subs) != 1 {
		t.Fatalf("Expected 1 submission, got %d", len(subs))
	}

	if subs[0].Name1 != "Romeo" || subs[0].Name2 != "Juliet" || subs[0].Percentage != 99 {
		t.Errorf("Submission data mismatch: %+v", subs[0])
	}
}
