package printline

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

// print line ExecutePrintLine("=", 1, "y")
func ExecutePrintLine(PrintValue string, blankRowValue string) {
	var ConsoleWide int = GetConsoleWide()
	// print symbol
	PrintOneLine(PrintValue, ConsoleWide, blankRowValue)
}

// ExecuteCenter execute cobra center command func.
// demo ExecuteCenter("hello world!", "=", "y", "y")
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

// PrintOneLine print one line
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

// 计算字符串在终端中的显示宽度
func DisplayWidth(s string) int {
	return runewidth.StringWidth(s)
}
