package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"os"
	"strings"
)

/*
main func
*/

func main() {
	BondCobra()

}

/*
Cobra code start
*/

func BondCobra() {
	// root command
	var rootCmd = &cobra.Command{
		// 输入什么会执行这个命令
		Use: "printline",
		// 简单描述信息
		Short: "start generator,A CLI tool to print lines",
	}

	// version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of printline.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("printline v0.1.0")
		},
		// 参数不能多余两个 ExactArgs(int) 参数不为n则报错 MinimumNArgs(int) 最少几个参数
		Args: cobra.MaximumNArgs(2),
	}

	// line command
	var lineCmd = &cobra.Command{
		Use:   "line",
		Short: "Print a line of specified characters,default symbol is =",
		Run: func(cmd *cobra.Command, args []string) {
			ExecutePrintLine(cmd, args)
		},
		Args: cobra.MaximumNArgs(1),
	}

	//center command
	var centerCmd = &cobra.Command{
		Use:   "center",
		Short: "Print the specified string in the center,default symbol is =",
		Run: func(cmd *cobra.Command, args []string) {
			ExecuteCenter(cmd, args)
		},
		Args: cobra.MaximumNArgs(1),
	}

	// Completely centered
	var CompletelyCenterCmd = &cobra.Command{
		Use:   "completely-center",
		Short: "Completely centered printing string, default symbol is =",
		Run: func(cmd *cobra.Command, args []string) {
			CompletelyCenterPrinting(cmd, args)
		},
		Args: cobra.MaximumNArgs(1),
	}

	// flags
	var printValue string       // Specify printing symbols
	var printLines int          // How many lines to print
	var printInformation string // print length info
	var BlankRow string

	// bond command
	lineCmd.Flags().IntVarP(&printLines, "lines", "l", 1, "How many lines to print.")
	lineCmd.Flags().StringVarP(&BlankRow, "blank row", "b", "n", "Whether to print a blank line at the beginning and end. If you want to enable it, use 'y'.")

	centerCmd.Flags().StringVarP(&printValue, "symbol", "s", "", "specify printing symbols.")
	centerCmd.Flags().StringVarP(&printInformation, "print info", "p", "", "Print detailed length information, If it needs to be enabled, please specify 'y'")
	centerCmd.Flags().StringVarP(&BlankRow, "blank row", "b", "n", "Whether to print a blank line at the beginning and end. If you want to enable it, use 'y'.")
	//centerCmd.Flags().IntVarP(&printLines, "lines", "l", 1, "How many lines to print.")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(lineCmd)
	rootCmd.AddCommand(CompletelyCenterCmd)
	rootCmd.AddCommand(centerCmd)

	// default :run rootCmd.Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
Cobra command execute func
*/

// ExecutePrintLine print line command func
func ExecutePrintLine(cmd *cobra.Command, args []string) {
	var PrintValue string = "="
	var ConsoleWide int

	// get args
	lineValue, err := cmd.Flags().GetInt("lines")
	blankRowValue, err := cmd.Flags().GetString("blank row")
	if err != nil {
		fmt.Printf("Failed to obtain parameters!\n%s", err)
		os.Exit(1)
	}

	if len(args) != 0 {
		if len(args[0]) != 1 {
			cmd.Println("Please enter a single character.")
			os.Exit(1)
		}
		PrintValue = args[0]
	}

	ConsoleWide = GetConsoleWide()
	// print symbol
	for i := 0; i < lineValue; i++ {
		PrintOneLine(PrintValue, ConsoleWide, blankRowValue)
	}
}

// ExecuteCenter execute cobra center command func.
func ExecuteCenter(cmd *cobra.Command, args []string) {
	// get cobra args
	strValue, err := cmd.Flags().GetString("symbol")
	PrintInfoValue, err := cmd.Flags().GetString("print info")
	blankRowValue, err := cmd.Flags().GetString("blank row")
	if err != nil {
		fmt.Printf("Failed to obtain parameters!\n%s", err)
	}

	// set print symbol
	var symbol string = "="
	if strValue != "" {
		symbol = strValue
	}

	// get console wide
	ConsoleWide := GetConsoleWide()

	// calculate console wide, string center
	if len(args) == 0 {
		fmt.Println("Please enter the characters that need to be printed!")
		os.Exit(1)
	}
	textRunes := []rune(args[0])
	textLength := len(textRunes)
	var SymbolWide int

	// print
	if textLength > ConsoleWide {
		PrintOneLine(symbol, ConsoleWide, blankRowValue)
		cmd.Println(args[0])
		PrintOneLine(symbol, ConsoleWide, blankRowValue)
	} else {
		SymbolLength := ConsoleWide - textLength - 2
		SymbolWide = SymbolLength / 2
		PrintCenter(symbol, SymbolLength, SymbolWide, args[0], blankRowValue)
	}

	// print info
	if PrintInfoValue == "y" {
		cmd.Printf("Console Length: %d, String Length: %d, Symbol Length: %d\n",
			ConsoleWide, textLength, SymbolWide)
	}
}

// CompletelyCenterPrinting :Completely Center Printing
func CompletelyCenterPrinting(cmd *cobra.Command, args []string) {
	fmt.Println("function not implemented ...")
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
	// 获取终端宽度
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to obtain terminal size: %v\n", err)
		os.Exit(1)
	}

	return width
}
