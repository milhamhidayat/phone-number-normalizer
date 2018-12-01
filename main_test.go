package main

import "testing"

type normalizeTestCase struct {
	input    string
	expected string
}

func TestNormalizer(t *testing.T) {
	testCases := []normalizeTestCase{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"123-456-7890", "1234567890"},
		{"1234567892", "1234567892"},
		{"(123) 456-7890", "1234567890"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := normalize(tc.input)
			if actual != tc.expected {
				t.Errorf("Got %s, but expected : %s", actual, tc.expected)
			}
		})
	}
}
