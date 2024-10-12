package editor

import (
	"testing"
)

// TestCountLeadingDashes tests the countLeadingDashes function.
func TestCountLeadingDashes(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"---file.go", 3},
		{"- dir", 1},
		{"--nested", 2},
		{"no-dashes", 0},
		{" -mixed", 1},
		{"", 0},
		{"\t-tab", 1},
	}

	for _, test := range tests {
		result := countLeadingDashes(test.input)
		if result != test.expected {
			t.Errorf("countLeadingDashes(%q) = %d; want %d", test.input, result, test.expected)
		}
	}
}

// TestIsFileLine tests the isFileLine function.
func TestIsFileLine(t *testing.T) {
	tests := []struct {
		input        string
		expectedIs   bool
		expectedName string
	}{
		{"-main.go", true, "main.go"},
		{"--utils.go", true, "utils.go"},
		{"-README:file", true, "README"},
		{"--LICENSE:file", true, "LICENSE"},
		{"-docs", false, "docs"},
		{"--src", false, "src"},
		{"-config", false, "config"},
		{"-script.sh", true, "script.sh"},
		{"-noextension:file", true, "noextension"},
		{"-invalid:fileextra", false, "invalid:fileextra"}, // Edge case
		{"", false, ""},
		{"---", false, ""},
	}

	for _, test := range tests {
		isFile, name := isFileLine(test.input)
		if isFile != test.expectedIs || name != test.expectedName {
			t.Errorf("isFileLine(%q) = (%v, %q); want (%v, %q)", test.input, isFile, name, test.expectedIs, test.expectedName)
		}
	}
}

// TestIsLineIncomplete tests the isLineIncomplete method.
func TestIsLineIncomplete(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"Complete line with content", "----file.go", false},
		{"Incomplete line with only dashes", "---", true},
		{"Incomplete line with dashes and spaces", "---   ", true},
		{"Complete line with mixed content", "---main.go", false},
		{"Empty line", "", true},
		{"Line with spaces only", "   ", true},
		{"Line with tabs and dashes", "\t--file.go", false},
	}

	for _, test := range tests {
		editor := NewEditor(nil) // statusBar is not needed for this test
		editor.Lines = []string{test.line}
		result := editor.isLineIncomplete(editor.Lines[0])
		if result != test.expected {
			t.Errorf("isLineIncomplete(%q) = %v; want %v", test.line, result, test.expected)
		}
	}
}

// TestEnforceDepth tests the enforceDepth method.
func TestEnforceDepth(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		currentY int
		line     string
		expected string
	}{
		{
			name:     "No depth enforcement needed",
			lines:    []string{"-src", "--main.go"},
			currentY: 1,
			line:     "--main.go",
			expected: "--main.go",
		},
		{
			name:     "Exceeds max depth",
			lines:    []string{"-src", "--main.go"},
			currentY: 1,
			line:     "---extra.go",
			expected: "--extra.go",
		},
		{
			name:     "Reduce depth",
			lines:    []string{"-src", "--main.go"},
			currentY: 1,
			line:     "----utils.go",
			expected: "--utils.go",
		},
		{
			name:     "Adjust depth based on parent",
			lines:    []string{"-src", "--main.go", "-docs"},
			currentY: 2,
			line:     "-docs",
			expected: "-docs",
		},
		{
			name:     "Initial line with no dashes",
			lines:    []string{""},
			currentY: 0,
			line:     "src",
			expected: "src",
		},
	}

	for _, test := range tests {
		editor := NewEditor(nil) // statusBar is not needed for this test
		editor.Lines = test.lines
		editor.cursorY = test.currentY
		result := editor.enforceDepth(test.line)
		if result != test.expected {
			t.Errorf("enforceDepth(%q) = %q; want %q", test.line, result, test.expected)
		}
	}
}

// TestValidateStructure tests the ValidateStructure method.
func TestValidateStructure(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		hasError bool
		errorMsg string
	}{
		{
			name:     "Valid structure",
			lines:    []string{"-src", "--main.go", "-docs", "--README.md"},
			hasError: false,
		},
		{
			name:     "Incomplete line",
			lines:    []string{"-src", "--", "-docs"},
			hasError: true,
			errorMsg: "line 2 is incomplete",
		},
		{
			name:     "Empty lines",
			lines:    []string{"-src", "", "-docs"},
			hasError: true,
			errorMsg: "line 2 is incomplete",
		},
		{
			name:     "All complete lines",
			lines:    []string{"-src", "--main.go", "--utils.go", "-docs", "--README.md"},
			hasError: false,
		},
		{
			name:     "Multiple incomplete lines",
			lines:    []string{"-", "--", "---file.go"},
			hasError: true,
			errorMsg: "line 1 is incomplete",
		},
		{
			name:     "No lines",
			lines:    []string{},
			hasError: false, // Depending on desired behavior; assuming no error
		},
	}

	for _, test := range tests {
		editor := NewEditor(nil)
		editor.Lines = test.lines
		err := editor.ValidateStructure()
		if test.hasError {
			if err == nil {
				t.Errorf("ValidateStructure(%q) expected error but got none", test.name)
			} else if err.Error() != test.errorMsg {
				t.Errorf("ValidateStructure(%q) error = %q; want %q", test.name, err.Error(), test.errorMsg)
			}
		} else {
			if err != nil {
				t.Errorf("ValidateStructure(%q) unexpected error: %v", test.name, err)
			}
		}
	}
}
