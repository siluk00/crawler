package main

import (
	"fmt"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectedErr bool
	}{
		{
			input:       "http://github.com/siluk00",
			expected:    "github.com/siluk00",
			expectedErr: false,
		}, {
			input:       "http://github.com/siluk00/",
			expected:    "github.com/siluk00",
			expectedErr: false,
		}, {
			input:       "https://github.com/siluk00/",
			expected:    "github.com/siluk00",
			expectedErr: false,
		}, {
			input:       "https://githubcom/siluk00/",
			expected:    "",
			expectedErr: true,
		}, {
			input:       "https:/github.com/siluk00/",
			expected:    "",
			expectedErr: true,
		}, {
			input:       "ht://github.com/siluk00/",
			expected:    "",
			expectedErr: true,
		}, {
			input:       "",
			expected:    "",
			expectedErr: true,
		}, {
			input:       "http://github.com/?commit=true",
			expected:    "github.com/?commit=true",
			expectedErr: false,
		}, {
			input:       "http://github.com/?commit=true&bring=sum+coif",
			expected:    "github.com/?bring=sum coif&commit=true",
			expectedErr: false,
		}, {
			input:       "http://github.com/?bring=sum+coif&commit=true",
			expected:    "github.com/?bring=sum coif&commit=true",
			expectedErr: false,
		},
	}

	for i, test := range tests {
		name := fmt.Sprintf("test %d", i)
		t.Run(name, func(t *testing.T) {
			actual, err := NormalizeURL(test.input)
			if err != nil {
				if !test.expectedErr {
					t.Errorf("Test %s failed: %v", name, err)
				}
				return
			}
			if actual != test.expected {
				t.Errorf("Test %s failed, input: %s, expected output: %s, actual: %s", name, test.input, test.expected, actual)
			}
		})
	}
}
