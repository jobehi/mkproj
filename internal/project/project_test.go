package project

import (
	"os"
	"path/filepath"
	"testing"
)

// TestBuildProjectStructure_Basic tests the creation of a basic project structure.
func TestBuildProjectStructure_Basic(t *testing.T) {
	rootDir := setupTestRootDir(t)
	defer os.RemoveAll(rootDir) // Clean up after the test

	lines := []string{
		"src",
		"-main.go",
		"-utils.go",
		"README.md",
	}

	// Call BuildProjectStructure
	BuildProjectStructure(lines, rootDir)

	// Validate the expected directory structure
	expectedDirs := []string{
		filepath.Join(rootDir, "src"),
	}
	expectedFiles := []string{
		filepath.Join(rootDir, "src", "main.go"),
		filepath.Join(rootDir, "src", "utils.go"),
		filepath.Join(rootDir, "README.md"),
	}

	validateStructure(t, expectedDirs, expectedFiles)
}

// TestBuildProjectStructure_FileWithoutExtension tests the creation of a file without an extension.
func TestBuildProjectStructure_FileWithoutExtension(t *testing.T) {
	rootDir := setupTestRootDir(t)
	defer os.RemoveAll(rootDir) // Clean up after the test

	lines := []string{
		"docs",
		"-README:file",
		"-LICENSE:file",
	}

	// Call BuildProjectStructure
	BuildProjectStructure(lines, rootDir)

	// Validate the expected directory structure
	expectedDirs := []string{
		filepath.Join(rootDir, "docs"),
	}
	expectedFiles := []string{
		filepath.Join(rootDir, "docs", "README"),
		filepath.Join(rootDir, "docs", "LICENSE"),
	}

	validateStructure(t, expectedDirs, expectedFiles)
}

// TestBuildProjectStructure_EmptyInput tests the behavior when given an empty input.
func TestBuildProjectStructure_EmptyInput(t *testing.T) {
	rootDir := setupTestRootDir(t)
	defer os.RemoveAll(rootDir) // Clean up after the test

	lines := []string{}

	// Call BuildProjectStructure
	BuildProjectStructure(lines, rootDir)

	// Ensure no directories or files were created
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		t.Errorf("Expected root directory %s to exist", rootDir)
	}

	files, err := os.ReadDir(rootDir)
	if err != nil {
		t.Fatalf("Error reading root directory: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("Expected no files or directories in root directory, but found %d", len(files))
	}
}

// Helper function to set up the test root directory
func setupTestRootDir(t *testing.T) string {
	rootDir, err := os.MkdirTemp("", "project_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	return rootDir
}

// Helper function to validate the created directories and files
func validateStructure(t *testing.T, expectedDirs []string, expectedFiles []string) {
	// Check directories
	for _, dir := range expectedDirs {
		if stat, err := os.Stat(dir); os.IsNotExist(err) || !stat.IsDir() {
			t.Errorf("Expected directory %s to exist", dir)
		}
	}

	// Check files
	for _, file := range expectedFiles {
		if stat, err := os.Stat(file); os.IsNotExist(err) || stat.IsDir() {
			t.Errorf("Expected file %s to exist", file)
		}
	}
}
