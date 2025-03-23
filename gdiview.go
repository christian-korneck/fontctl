package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tadvi/winc"
)

const (
	fontStyleRegular byte = 0
	fontStyleBold    byte = 1 << (iota - 1)
	fontStyleItalic
	fontStyleUnderline
	fontStyleStrikeOut
)

func PreviewFontWithGDI(fontName string, fontStyle string) {

	var fontStyleB byte

	switch strings.ToLower(fontStyle) {
	case "regular":
		fontStyleB = fontStyleRegular
	case "bold":
		fontStyleB = fontStyleBold
	case "bold-italic":
		fontStyleB = fontStyleBold | fontStyleItalic
	case "italic":
		fontStyleB = fontStyleItalic
	default:
		// as this is a separate gui app, it's ok to let it pollute stdout
		fmt.Fprintf(os.Stderr, "Warning - invalid font style %s, falling back to 'regular'.\n", fontStyle)
		fontStyleB = fontStyleRegular
	}

	mainWindow := winc.NewForm(nil)
	mainWindow.SetText("fontctl Preview")
	mainWindow.SetSize(750, 400)
	mainWindow.SetMinSize(750, 400)
	mainWindow.SetMaxSize(750, 400)

	menu := mainWindow.NewMenu()
	fileMn := menu.AddSubMenu("File")
	exit := fileMn.AddItem("Exit", winc.Shortcut{Modifiers: winc.ModAlt, Key: winc.KeyF4})
	exit.OnClick().Bind(func(e *winc.Event) {
		winc.Exit()
	})
	menu.Show()

	font := winc.NewFont(fontName, 18, fontStyleB) // fixed font size 18pt
	textBox := winc.NewMultiEdit(mainWindow)
	textBox.SetFont(font)
	textBox.SetMinSize(720, 300)
	textBox.SetMaxSize(720, 300)

	textBox.SetText("The quick brown fox jumps over the lazy dog. 1234567890")
	textBox.AddLine("abcdefghijklmnopqrstuvwxyz")
	textBox.AddLine("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	winc.RunMainLoop() // start gui loop
}

func wndOnClose(arg *winc.Event) {
	winc.Exit()
}

func DemoWithGDI() {
	PreviewFontWithGDI("Comic Sans MS", "bold-italic")
}
