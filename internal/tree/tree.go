package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DisplayDirectoryTree shows the directory tree.
func DisplayDirectoryTree(rootDir string, showHidden bool) {
	fmt.Println("Current Directory Structure:")
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err
		}
		if path == rootDir {
			return nil
		}
		relativePath, _ := filepath.Rel(rootDir, path)
		if !showHidden {
			// Skip hidden files and directories
			parts := strings.Split(relativePath, string(os.PathSeparator))
			for _, part := range parts {
				if strings.HasPrefix(part, ".") {
					if info.IsDir() && info.Name() == part {
						return filepath.SkipDir
					}
					return nil
				}
			}
		}
		depth := strings.Count(relativePath, string(os.PathSeparator))
		indent := strings.Repeat("-", depth)
		name := info.Name()
		if info.IsDir() {
			fmt.Printf("%s%s\n", indent, name)
		} else {
			if filepath.Ext(name) == "" {
				fmt.Printf("%s%s:file\n", indent, name)
			} else {
				fmt.Printf("%s%s\n", indent, name)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error displaying directory tree: %v\n", err)
	}
}
