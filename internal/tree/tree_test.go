package tree

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// Utility function to capture stdout for testing purposes.
func captureOutput(f func()) string {
	var buf bytes.Buffer
	// Save the original stdout
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()

	// Temporarily redirect stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the function to capture output
	f()

	// Restore stdout and get the output
	w.Close()
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

// TestDisplayDirectoryTree_EmptyDirectory tests that an empty directory outputs nothing after "Current Directory Structure:"
func TestDisplayDirectoryTree_EmptyDirectory(t *testing.T) {
	// Setup: Create a temporary empty directory
	rootDir := setupEmptyTestDirectory(t)
	defer os.RemoveAll(rootDir) // Cleanup after test

	// Capture the output of the DisplayDirectoryTree function
	output := captureOutput(func() {
		DisplayDirectoryTree(rootDir, false)
	})

	// Expected directory structure output for an empty directory
	expected := "Current Directory Structure:"

	if strings.TrimSpace(output) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

// Helper function to create an empty temporary directory for testing
func setupEmptyTestDirectory(t *testing.T) string {
	rootDir, err := os.MkdirTemp("", "emptydir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return rootDir
}
