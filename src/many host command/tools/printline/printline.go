package printline

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

// ExecutePrintLine prints a separator line across the console width.
// PrintValue: The character to use for the line (e.g., "=").
// blankRowValue: If "y", adds a blank line before and after the separator.
func ExecutePrintLine(PrintValue string, blankRowValue string) {
	var ConsoleWide int = GetConsoleWide()
	// print symbol
	PrintOneLine(PrintValue, ConsoleWide, blankRowValue)
}

// ExecuteCenter prints text centered in the console, surrounded by separator characters.
//
// content: The text to display.
// symbol: The character to use for padding (e.g., "=").
// blankRowValue: If "y", adds a blank line before and after the output.
// PrintInfoValue: If "y", prints debug info about lengths.
func ExecuteCenter(content string, symbol string, blankRowValue string, PrintInfoValue string) {
	// get console wide
	ConsoleWide := GetConsoleWide()
	// textRunes := []rune(content)
	// textLength := len(textRunes)
	textLength := DisplayWidth(content)

	var SymbolWide int

	// print
	if textLength > ConsoleWide {
		PrintOneLine(symbol, ConsoleWide, blankRowValue)
		fmt.Println(content)
		PrintOneLine(symbol, ConsoleWide, blankRowValue)
	} else {
		SymbolLength := ConsoleWide - textLength - 2
		SymbolWide = SymbolLength / 2
		PrintCenter(symbol, SymbolLength, SymbolWide, content, blankRowValue)
	}

	// print info
	if PrintInfoValue == "y" {
		fmt.Printf("Console Length: %d, String Length: %d, Symbol Length: %d\n",
			ConsoleWide, textLength, SymbolWide)
	}
}

/*
Print func
*/

// PrintOneLine prints a full line of the specified character.
//
// PrintValue: The character to repeat.
// ConsoleWide: The width of the line.
// blankRowValue: If "y", adds blank lines around the output.
func PrintOneLine(PrintValue string, ConsoleWide int, blankRowValue string) {
	line := strings.Repeat(PrintValue, ConsoleWide)

	if blankRowValue == "y" {
		fmt.Println()
		fmt.Println(line)
		fmt.Println()
	} else {
		fmt.Println(line)
	}
}

// PrintCenter formats and prints the content centered with padding symbols.
func PrintCenter(PrintValue string, SymbolLength int, SymbolWide int, PrintStr string, blankRowValue string) {
	SymbolLine := strings.Repeat(PrintValue, SymbolWide)
	// judge: Print a blank line at the beginning
	if blankRowValue == "y" {
		fmt.Println()
	}

	if SymbolLength%2 == 0 {
		// There is no remainder, the length is just right
		fmt.Println(SymbolLine + " " + PrintStr + " " + SymbolLine)
	} else {
		// There is a remainder, and there will be 1 empty space after it
		fmt.Println(SymbolLine + " " + PrintStr + " " + SymbolLine + PrintValue)
	}

	// Determine whether to print a blank line at the end
	if blankRowValue == "y" {
		fmt.Println()
	}
}

/*
Console related code
*/

// GetConsoleWide returns the width of the current terminal.
// If it's not a terminal (e.g. pipe), it returns a default width of 40.
func GetConsoleWide() int {
	// 检查是否存在伪终端 也就是管道符 default使用20宽度
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return 40
	}

	// 获取终端宽度
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to obtain terminal size: %v\n", err)
		os.Exit(1)
	}

	return width
}

// DisplayWidth calculates the visual width of a string in the terminal.
// It handles multi-byte characters correctly.
func DisplayWidth(s string) int {
	return runewidth.StringWidth(s)
}
