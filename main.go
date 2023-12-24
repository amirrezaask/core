/*
   move all functionalities to editor struct
   editor.GetCursorBufferIndex()
*/

package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	font     rl.Font
	fontSize float32
)

type BufferID int64
type Buffer struct {
	Cursor    Position
	Content   []byte
	FilePath  string
	Keymaps   []Keymap
	Variables Variables
	Commands  Commands
}

type Editor struct {
	Buffers           []Buffer
	GlobalKeymaps     []Keymap
	GlobalVariables   Variables
	Commands          Commands
	Windows           []Window
	ActiveWindowIndex int
	LineWrapping      bool
}

func (e Editor) CurrentBuffer() *Buffer {
	return &e.Buffers[e.Windows[e.ActiveWindowIndex].BufferIndex]
}
func (e Editor) CurrentWindow() *Window {
	return &e.Windows[e.ActiveWindowIndex]
}

type Command func(Editor) error
type Variables map[string]any
type Key struct {
	Ctrl  bool
	Alt   bool
	Shift bool
	Super bool
	K     string
}
type Keymap map[Key]Command
type Commands map[string]Command
type Position struct {
	Line   int
	Column int
}

func (p Position) String() string {
	return fmt.Sprintf("Line: %d Column:%d\n", p.Line, p.Column)
}



type WindowID int64
type Window struct {
	BufferIndex  int
	zeroLocation rl.Vector2
	Height       int
	Width        int
	Cursor       Position
	VisualLines  []visualLine
}


type visualLine struct {
	visualLineIndex int
	startIndex      int
	endIndex        int
	ActualLine      int
}

// we are considering fonts to mono spaced,
func (e *Editor) RenderBufferInWindow(buffer *Buffer, window *Window) {
	//first scan through buffer.Contents
	// every new line adds a visual line
	// every time we reach windowMaxColumn we add visualLine
	window.VisualLines = []visualLine{}
	charSize := measureTextSize(font, 'a', fontSize, 0)
	totalVisualLines := 0
	lineCharCounter := 0
	var actualLineIndex int
	var start int
	windowMaxColumn := window.Width  / int(charSize.X)
	windowMaxLine := window.Height / int (charSize.Y)
	for idx, char := range buffer.Content {
		lineCharCounter++
		if char == '\n' {
			fmt.Printf("At %d, saw a new line char\n", idx)
			window.VisualLines = append(window.VisualLines, visualLine{
				visualLineIndex: totalVisualLines,
				startIndex: start,
				endIndex: idx,
				ActualLine: actualLineIndex,
			})
			totalVisualLines ++
			actualLineIndex++
			lineCharCounter = 0
			start = idx+1
		}

		if lineCharCounter > windowMaxColumn {
			fmt.Printf("At %d, need to break line\n", idx)

			window.VisualLines = append(window.VisualLines, visualLine{
				visualLineIndex: totalVisualLines,
				startIndex: start,
				endIndex: idx,
				ActualLine: actualLineIndex,
				
			})
			totalVisualLines++
			lineCharCounter=0
			start = idx
		}
		

	}

	for _, line := range window.VisualLines {
		fmt.Printf("%+v\n", line)
		fmt.Printf("Y %f\n", float32(line.visualLineIndex)*charSize.Y)
		fmt.Println(string(buffer.Content[line.startIndex:line.endIndex]))
		if line.visualLineIndex > windowMaxLine {
			break
		}
		rl.DrawTextEx(font,
			string(buffer.Content[line.startIndex:line.endIndex]),
			rl.Vector2{X: window.zeroLocation.X, Y: float32(line.visualLineIndex)*charSize.Y},
			fontSize,
			0,
			rl.White)
	}

	
}


func (buffer *Buffer) InsertCharAtCursor(char byte) error {
	// buffer.Content[buffer.Cursor.Line] = append(buffer.Content[buffer.Cursor.Line][0:buffer.Cursor.Column+1], buffer.Content[buffer.Cursor.Line][buffer.Cursor.Column:]...)
	// buffer.Content[buffer.Cursor.Line][buffer.Cursor.Column] = char
	// buffer.Cursor.Column = buffer.Cursor.Column + 1

	// return nil
	return nil
}

func main() {
	// basic setup
	rl.InitWindow(1920, 1080, "core editor")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// create editor
	editor := Editor{
		LineWrapping: true,
	}

	fontSize = 20
	rl.SetTextLineSpacing(int(fontSize))
	rl.SetMouseCursor(rl.MouseCursorIBeam)
	editor.Buffers = append(editor.Buffers, Buffer{
		Cursor:   Position{0, 0},
		Content:  []byte(`orem ipsum dolor sit amet . The graphic and typographic operators know this well, in reality all the professions dealing with the universe of communication have a stable relationship with these words, but what is it? Lorem ipsum is a dummy text without any sense.

It is a sequence of Latin words that, as they are positioned, do not form sentences with a complete sense, but give life to a test text useful to fill spaces that will subsequently be occupied from ad hoc texts composed by communication professionals.

It is certainly the most famous placeholder text even if there are different versions distinguishable from the order in which the Latin words are repeated.

Lorem ipsum contains the typefaces more in use, an aspect that allows you to have an overview of the rendering of the text in terms of font choice and font size .

When referring to Lorem ipsum, different expressions are used, namely fill text , fictitious text , blind text or placeholder text : in short, its meaning can also be zero, but its usefulness is so clear as to go through the centuries and resist the ironic and modern versions that came with the arrival of the web.l`),
		FilePath: "test.txt",
	})
	editor.Windows = append(editor.Windows, Window{
		BufferIndex: 0,
		zeroLocation: rl.Vector2{
			X: 0, Y: 0,
		},
		Height: rl.GetRenderHeight(),
		Width:  rl.GetRenderWidth(),
	})

	font = rl.LoadFontEx("FiraCode.ttf", int32(fontSize), nil)
	for !rl.WindowShouldClose() {
		buffer := &editor.Buffers[editor.Windows[editor.ActiveWindowIndex].BufferIndex]

		// execute any command that should be executed
		cmd := defaultKeymap[MakeKey(buffer)]
		if cmd != nil {
			if err := cmd(editor); err != nil {
				panic(err)
			}
		}

		// Render
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		editor.RenderBufferInWindow(&editor.Buffers[0], &editor.Windows[0])

		rl.EndDrawing()
	}

}

var charSizeCache = map[byte]rl.Vector2{} //TODO: if font size or font changes this is fucked
func measureTextSize(font rl.Font, s byte, size float32, spacing float32) rl.Vector2 {
	if charSize, exists := charSizeCache[s]; exists {
		return charSize
	}
	charSize := rl.MeasureTextEx(font, string(s), size, spacing)
	charSizeCache[s] = charSize
	return charSize
}
