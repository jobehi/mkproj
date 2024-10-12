package editor

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Editor struct {
	*tview.Box
	Lines            []string
	cursorX, cursorY int
	statusBar        *tview.TextView
}

// NewEditor creates a new Editor instance.
func NewEditor(statusBar *tview.TextView) *Editor {
	return &Editor{
		Box:       tview.NewBox(),
		Lines:     []string{""},
		statusBar: statusBar,
	}
}

// Draw renders the editor on the screen.
func (e *Editor) Draw(screen tcell.Screen) {
	e.Box.DrawForSubclass(screen, e)
	defStyle := tcell.StyleDefault
	x, y, width, height := e.GetInnerRect()
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			screen.SetContent(x+col, y+row, ' ', nil, defStyle)
		}
	}
	startLine := 0
	if len(e.Lines) > height {
		startLine = len(e.Lines) - height
	}
	for i := startLine; i < len(e.Lines); i++ {
		line := e.Lines[i]
		tview.Print(screen, line, x, y+i-startLine, width, tview.AlignLeft, tcell.ColorWhite)
	}
	screen.ShowCursor(x+e.cursorX, y+e.cursorY)
}

// InputHandler handles key events for the editor.
func (e *Editor) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return e.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		line := e.Lines[e.cursorY]
		switch event.Key() {
		case tcell.KeyTab:
			if e.cursorX > len(line) {
				e.cursorX = len(line)
			}
			line = line[:e.cursorX] + "-" + line[e.cursorX:]
			e.cursorX++
			e.Lines[e.cursorY] = e.enforceDepth(line)
		case tcell.KeyRune:
			ch := event.Rune()
			if ch == '\t' {
				ch = '-'
			} else if ch == ' ' {
				return
			}
			if e.cursorX > len(line) {
				e.cursorX = len(line)
			}
			line = line[:e.cursorX] + string(ch) + line[e.cursorX:]
			e.cursorX++
			e.Lines[e.cursorY] = e.enforceDepth(line)
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if e.cursorX > len(line) {
				e.cursorX = len(line)
			}
			if e.cursorX > 0 {
				line = line[:e.cursorX-1] + line[e.cursorX:]
				e.Lines[e.cursorY] = line
				e.cursorX--
			} else if e.cursorY > 0 {
				prevLine := e.Lines[e.cursorY-1]
				e.cursorX = len(prevLine)
				e.Lines[e.cursorY-1] = prevLine + line
				e.Lines = append(e.Lines[:e.cursorY], e.Lines[e.cursorY+1:]...)
				e.cursorY--
			}
		case tcell.KeyDelete:
			if e.cursorX < len(line) {
				line = line[:e.cursorX] + line[e.cursorX+1:]
				e.Lines[e.cursorY] = line
			} else if e.cursorY < len(e.Lines)-1 {
				nextLine := e.Lines[e.cursorY+1]
				e.Lines[e.cursorY] = line + nextLine
				e.Lines = append(e.Lines[:e.cursorY+1], e.Lines[e.cursorY+2:]...)
			}
		case tcell.KeyLeft:
			if e.cursorX > 0 {
				e.cursorX--
			} else if e.cursorY > 0 {
				e.cursorY--
				e.cursorX = len(e.Lines[e.cursorY])
			}
		case tcell.KeyRight:
			if e.cursorX < len(line) {
				e.cursorX++
			} else if e.cursorY < len(e.Lines)-1 {
				e.cursorY++
				e.cursorX = 0
			}
		case tcell.KeyUp:
			if e.cursorY > 0 {
				e.cursorY--
				if e.cursorX > len(e.Lines[e.cursorY]) {
					e.cursorX = len(e.Lines[e.cursorY])
				}
			}
		case tcell.KeyDown:
			if e.cursorY < len(e.Lines)-1 {
				e.cursorY++
				if e.cursorX > len(e.Lines[e.cursorY]) {
					e.cursorX = len(e.Lines[e.cursorY])
				}
			}
		case tcell.KeyEnter:
			if e.cursorX > len(line) {
				e.cursorX = len(line)
			}
			if e.isLineIncomplete(e.Lines[e.cursorY]) {
				e.statusBar.SetText("Cannot add a new line after an incomplete line.")
				return
			}
			newLine := e.Lines[e.cursorY][e.cursorX:]
			e.Lines[e.cursorY] = e.Lines[e.cursorY][:e.cursorX]
			e.Lines = append(e.Lines[:e.cursorY+1], append([]string{newLine}, e.Lines[e.cursorY+1:]...)...)
			e.cursorY++
			e.cursorX = 0
		}
		e.Lines[e.cursorY] = e.enforceDepth(e.Lines[e.cursorY])
		e.statusBar.SetText("")
	})
}

// enforceDepth enforces depth restrictions.
func (e *Editor) enforceDepth(line string) string {
	line = strings.TrimLeft(line, " ")
	line = strings.ReplaceAll(line, "\t", "-")
	maxDepth := e.getMaxAllowedDepth(e.cursorY)
	if maxDepth < 0 {
		maxDepth = 0
	}
	dashCount := countLeadingDashes(line)
	if dashCount > maxDepth {
		line = strings.TrimLeft(line, "-")
		line = strings.Repeat("-", maxDepth) + line
	}
	return line
}

// getMaxAllowedDepth returns the maximum allowed depth.
func (e *Editor) getMaxAllowedDepth(currentIndex int) int {
	if currentIndex == 0 {
		return 0
	}
	for i := currentIndex - 1; i >= 0; i-- {
		lineContent := e.Lines[i]
		if lineContent == "" {
			continue
		}
		depth := countLeadingDashes(lineContent)
		isFile, _ := isFileLine(lineContent)
		if !isFile {
			return depth + 1
		} else {
			return depth
		}
	}
	return 0
}

// isLineIncomplete checks if a line is incomplete.
func (e *Editor) isLineIncomplete(line string) bool {
	trimmed := strings.TrimLeft(line, "-")
	return strings.TrimSpace(trimmed) == ""
}

// ValidateStructure checks if the structure is valid.
func (e *Editor) ValidateStructure() error {
	for i, line := range e.Lines {
		if e.isLineIncomplete(line) {
			return fmt.Errorf("line %d is incomplete", i+1)
		}
	}
	return nil
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
