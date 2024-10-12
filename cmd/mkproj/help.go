package main

import "fmt"

func printHelp() {
	fmt.Println(`mkproj: A Simple CLI Tool to Grow Your Project Trees 🌳

Usage:
  mkproj [command] [options]

Commands:
  create       Create a project structure from a text file or piped input
  tree         Display the current directory structure
  help         Display this help message

Options:
  --root=<path>    Specify the root directory for your project structure (default is current directory)
  --file=<path>    Provide a file that contains the project structure (used with 'create')

Interactive Mode:
  By default, mkproj starts in interactive mode where you can manually build your project structure.
  Use standard editing keys to modify the structure.
  Press F2 to save and create the structure.
  Press Esc to exit without saving.

Examples:
  # Start mkproj in interactive mode
  mkproj

  # Create a project structure from a text file
  mkproj create --file=structure.txt --root=./new_project

  # Display the current directory tree without hidden files
  mkproj tree --root=./my_project

  # Display the current directory tree including hidden files
  mkproj tree --root=./my_project --all`)
}
