package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/jobehi/mkproj/internal/editor"
	"github.com/jobehi/mkproj/internal/project"
	"github.com/jobehi/mkproj/internal/tree"
	"github.com/rivo/tview"
)

var rootDir string
var inputFile string

func main() {
	// Parse flags but not immediately
	rootFlag := flag.String("root", ".", "Root directory for project structure")
	fileFlag := flag.String("file", "", "Input file with project structure")
	flag.Usage = printHelp

	// Parse the command (e.g., "tree", "create", etc.)
	if len(os.Args) < 2 {
		runInteractiveMode(*rootFlag)
		return
	}

	command := os.Args[1]
	args := os.Args[2:]
	flag.CommandLine.Parse(args) // Parse the flags after the command

	rootDir = *rootFlag
	inputFile = *fileFlag

	// Handle help command
	if command == "help" {
		printHelp()
		return
	}

	// Handle tree command
	if command == "tree" {
		treeFlags := flag.NewFlagSet("tree", flag.ExitOnError)
		allFlag := treeFlags.Bool("all", false, "Include hidden files and directories")
		allFlagShort := treeFlags.Bool("a", false, "Include hidden files and directories (shorthand)")
		rootFlag := treeFlags.String("root", ".", "Root directory for project structure")
		treeFlags.Parse(args)

		rootDir = *rootFlag
		showHidden := *allFlag || *allFlagShort

		tree.DisplayDirectoryTree(rootDir, showHidden)
		return
	}

	// Handle create command
	if command == "create" {
		if inputFile != "" {
			// Read from input file
			file, err := os.Open(inputFile)
			if err != nil {
				fmt.Printf("Error reading input file %s: %v\n", inputFile, err)
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			var structure []string
			for scanner.Scan() {
				structure = append(structure, scanner.Text())
			}
			project.BuildProjectStructure(structure, rootDir)
			return
		}

		// Handle piped input
		if isPipedInput() {
			reader := bufio.NewReader(os.Stdin)
			var structure []string
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				structure = append(structure, strings.TrimSpace(line))
			}
			project.BuildProjectStructure(structure, rootDir)
			return
		}
	}

	// If no command matches, run the interactive mode by default
	runInteractiveMode(rootDir)
}

// isPipedInput detects if there is piped input from stdin
func isPipedInput() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice == 0
}

// runInteractiveMode launches the interactive mode for project structure building
func runInteractiveMode(rootDir string) {
	app := tview.NewApplication()

	// Status bar for feedback messages
	statusBar := tview.NewTextView().SetDynamicColors(true).SetText("Ready").SetTextAlign(tview.AlignLeft)

	// Instructions
	instructions := tview.NewTextView().
		SetText(fmt.Sprintf("Root Directory: %s\n", rootDir) +
			"Welcome to mkproj\n" +
			"Enter your project structure below.\n" +
			"Use tabs for depth and filename:file for files without extensions.\n" +
			"Press F2 to save and create the structure, Esc to quit.").
		SetDynamicColors(true)

	// Create the editor
	ed := editor.NewEditor(statusBar)

	// Capture inputs like F2 (save) and Esc (exit)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF2:
			if err := ed.ValidateStructure(); err != nil {
				statusBar.SetText(fmt.Sprintf("Validation error: %v", err))
				return nil
			}
			app.Stop()
			project.BuildProjectStructure(ed.Lines, rootDir)
			return nil
		case tcell.KeyEsc:
			app.Stop()
			return nil
		}
		return event
	})

	// Layout: instructions at top, editor in the middle, and status bar at the bottom
	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(instructions, 8, 1, false).
		AddItem(ed, 0, 1, true).
		AddItem(statusBar, 1, 1, false)

	// Run the interactive mode application
	if err := app.SetRoot(layout, true).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
