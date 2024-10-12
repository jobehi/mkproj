package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BuildProjectStructure builds the project structure from lines
func BuildProjectStructure(lines []string, rootDir string) {
	fmt.Println("Building project structure... Hold on tight! 🛠️")
	err := os.MkdirAll(rootDir, 0755)
	if err != nil {
		fmt.Printf("Error creating root directory %s: %v\n", rootDir, err)
		return
	}
	pathStack := []string{rootDir}
	for _, line := range lines {
		content := strings.TrimRight(line, "\r\n")
		if content == "" {
			continue
		}
		depth := countLeadingDashes(content)
		if depth < 0 {
			depth = 0
		}
		if depth > len(pathStack)-1 {
			depth = len(pathStack) - 1
		}
		pathStack = pathStack[:depth+1]
		parentDir := pathStack[len(pathStack)-1]
		isFile, name := isFileLine(content)
		if name == "" {
			fmt.Printf("Invalid name at line: %s\n", content)
			continue
		}
		fullPath := filepath.Join(parentDir, name)
		if isFile {
			file, err := os.Create(fullPath)
			if err != nil {
				fmt.Printf("Error creating file %s: %v\n", fullPath, err)
				continue
			}
			file.Close()
			fmt.Printf("Created file: %s\n", fullPath)
		} else {
			err := os.Mkdir(fullPath, 0755)
			if err != nil {
				fmt.Printf("Error creating directory %s: %v\n", fullPath, err)
				continue
			}
			fmt.Printf("Created directory: %s\n", fullPath)
			pathStack = append(pathStack, fullPath)
		}
	}
	displayFinalStructure(rootDir)
}

// displayFinalStructure shows the final structure.
func displayFinalStructure(rootDir string) {
	fmt.Println("\nFinal Project Structure:")
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err
		}
		relativePath, _ := filepath.Rel(rootDir, path)
		if relativePath != "." {
			depth := strings.Count(relativePath, string(os.PathSeparator))
			fmt.Printf("%s%s\n", strings.Repeat("  ", depth), info.Name())
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", rootDir, err)
	}
}

// countLeadingDashes counts the number of leading dashes.
func countLeadingDashes(s string) int {
	count := 0
	for _, char := range s {
		if char == '-' {
			count++
		} else if char == ' ' || char == '\t' {
			continue
		} else {
			break
		}
	}
	return count
}

// isFileLine checks if a line represents a file.
func isFileLine(line string) (bool, string) {
	name := strings.TrimLeft(line, "-")
	name = strings.TrimSpace(name)
	isFile := false
	if strings.HasSuffix(name, ":file") {
		isFile = true
		name = strings.TrimSuffix(name, ":file")
		name = strings.TrimSpace(name)
	} else if strings.Contains(name, ".") {
		isFile = true
	}
	return isFile, name
}
