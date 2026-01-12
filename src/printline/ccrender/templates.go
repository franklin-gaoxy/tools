package ccrender

import (
	"errors"
	"fmt"
	"strings"
)

type Config struct {
	TotalWidth   int
	Text         string
	BlankRow     bool
	LeftPadding  int
	RightPadding int
	Symbol       string
}

type Renderer func(cfg Config) ([]string, error)

var templates = map[string]Renderer{
	"box":     renderBox,
	"ascii":   renderASCII,
	"solid":   renderSolid,
	"double":  renderDouble,
	"log":     renderLog,
	"warn":    renderWarn,
	"error":   renderError,
	"success": renderSuccess,
}

func ComputeAvailableWidth(consoleWidth int, headerPrefix string) (int, bool) {
	available := consoleWidth - len([]rune(headerPrefix))
	if available <= 0 {
		return 0, false
	}
	return available, true
}

func Render(templateKey string, cfg Config) ([]string, error) {
	key := strings.ToLower(strings.TrimSpace(templateKey))
	if key == "" {
		key = "box"
	}

	renderer, ok := templates[key]
	if !ok {
		return nil, fmt.Errorf("unknown template/style %q", templateKey)
	}

	return renderer(cfg)
}

func normalizeSingleRuneSymbol(symbol string) (rune, bool) {
	if symbol == "" {
		symbol = "="
	}
	runes := []rune(symbol)
	if len(runes) != 1 {
		return 0, false
	}
	return runes[0], true
}

func renderBox(cfg Config) ([]string, error) {
	if cfg.TotalWidth < 4 {
		return nil, errors.New("terminal width is too small")
	}

	innerWidth := cfg.TotalWidth - 2
	top := "┌" + strings.Repeat("─", innerWidth) + "┐"
	bottom := "└" + strings.Repeat("─", innerWidth) + "┘"
	emptyInnerLine := "│" + strings.Repeat(" ", innerWidth) + "│"
	textLine := "│" + strings.Repeat(" ", cfg.LeftPadding) + cfg.Text + strings.Repeat(" ", cfg.RightPadding) + "│"

	lines := []string{top}
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, textLine)
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, bottom)

	return lines, nil
}

func renderASCII(cfg Config) ([]string, error) {
	if cfg.TotalWidth < 4 {
		return nil, errors.New("terminal width is too small")
	}

	innerWidth := cfg.TotalWidth - 2
	top := "+" + strings.Repeat("-", innerWidth) + "+"
	bottom := "+" + strings.Repeat("-", innerWidth) + "+"
	emptyInnerLine := "|" + strings.Repeat(" ", innerWidth) + "|"
	textLine := "|" + strings.Repeat(" ", cfg.LeftPadding) + cfg.Text + strings.Repeat(" ", cfg.RightPadding) + "|"

	lines := []string{top}
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, textLine)
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, bottom)

	return lines, nil
}

func renderSolid(cfg Config) ([]string, error) {
	if cfg.TotalWidth < 4 {
		return nil, errors.New("terminal width is too small")
	}

	symbolRune, ok := normalizeSingleRuneSymbol(cfg.Symbol)
	if !ok {
		return nil, errors.New("please specify a single character for --symbol")
	}

	innerWidth := cfg.TotalWidth - 2
	borderLine := strings.Repeat(string(symbolRune), cfg.TotalWidth)
	emptyInnerLine := string(symbolRune) + strings.Repeat(" ", innerWidth) + string(symbolRune)
	textLine := string(symbolRune) + strings.Repeat(" ", cfg.LeftPadding) + cfg.Text + strings.Repeat(" ", cfg.RightPadding) + string(symbolRune)

	lines := []string{borderLine}
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, textLine)
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, borderLine)

	return lines, nil
}

func renderDouble(cfg Config) ([]string, error) {
	if cfg.TotalWidth < 4 {
		return nil, errors.New("terminal width is too small")
	}

	innerWidth := cfg.TotalWidth - 2
	top := "╔" + strings.Repeat("═", innerWidth) + "╗"
	bottom := "╚" + strings.Repeat("═", innerWidth) + "╝"
	emptyInnerLine := "║" + strings.Repeat(" ", innerWidth) + "║"
	textLine := "║" + strings.Repeat(" ", cfg.LeftPadding) + cfg.Text + strings.Repeat(" ", cfg.RightPadding) + "║"

	lines := []string{top}
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, textLine)
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, bottom)

	return lines, nil
}

func renderLog(cfg Config) ([]string, error) {
	if cfg.TotalWidth < 8 {
		return nil, errors.New("terminal width is too small")
	}

	innerWidth := cfg.TotalWidth - 2
	label := " LOG "
	labelWidth := len([]rune(label))
	if innerWidth <= labelWidth {
		return renderASCII(cfg)
	}

	top := "+" + label + strings.Repeat("-", innerWidth-labelWidth) + "+"
	bottom := "+" + strings.Repeat("-", innerWidth) + "+"
	emptyInnerLine := "|" + strings.Repeat(" ", innerWidth) + "|"
	textLine := "|" + strings.Repeat(" ", cfg.LeftPadding) + cfg.Text + strings.Repeat(" ", cfg.RightPadding) + "|"

	lines := []string{top}
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, textLine)
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, bottom)

	return lines, nil
}

func renderWarn(cfg Config) ([]string, error) {
	if cfg.TotalWidth < 8 {
		return nil, errors.New("terminal width is too small")
	}

	innerWidth := cfg.TotalWidth - 2
	label := " WARN "
	labelWidth := len([]rune(label))
	if innerWidth <= labelWidth {
		return renderASCII(cfg)
	}

	top := "!" + label + strings.Repeat("!", innerWidth-labelWidth) + "!"
	bottom := strings.Repeat("!", cfg.TotalWidth)
	emptyInnerLine := "!" + strings.Repeat(" ", innerWidth) + "!"
	textLine := "!" + strings.Repeat(" ", cfg.LeftPadding) + cfg.Text + strings.Repeat(" ", cfg.RightPadding) + "!"

	lines := []string{top}
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, textLine)
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, bottom)

	return lines, nil
}

func renderError(cfg Config) ([]string, error) {
	if cfg.TotalWidth < 8 {
		return nil, errors.New("terminal width is too small")
	}

	innerWidth := cfg.TotalWidth - 2
	label := " ERROR "
	labelWidth := len([]rune(label))
	if innerWidth <= labelWidth {
		return renderSolid(Config{
			TotalWidth:   cfg.TotalWidth,
			Text:         cfg.Text,
			BlankRow:     cfg.BlankRow,
			LeftPadding:  cfg.LeftPadding,
			RightPadding: cfg.RightPadding,
			Symbol:       "#",
		})
	}

	top := "#" + label + strings.Repeat("#", innerWidth-labelWidth) + "#"
	bottom := strings.Repeat("#", cfg.TotalWidth)
	emptyInnerLine := "#" + strings.Repeat(" ", innerWidth) + "#"
	textLine := "#" + strings.Repeat(" ", cfg.LeftPadding) + cfg.Text + strings.Repeat(" ", cfg.RightPadding) + "#"

	lines := []string{top}
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, textLine)
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, bottom)

	return lines, nil
}

func renderSuccess(cfg Config) ([]string, error) {
	if cfg.TotalWidth < 8 {
		return nil, errors.New("terminal width is too small")
	}

	innerWidth := cfg.TotalWidth - 2
	label := " OK "
	labelWidth := len([]rune(label))
	if innerWidth <= labelWidth {
		return renderASCII(cfg)
	}

	top := "[" + label + strings.Repeat("=", innerWidth-labelWidth) + "]"
	bottom := "[" + strings.Repeat("=", innerWidth) + "]"
	emptyInnerLine := "[" + strings.Repeat(" ", innerWidth) + "]"
	textLine := "[" + strings.Repeat(" ", cfg.LeftPadding) + cfg.Text + strings.Repeat(" ", cfg.RightPadding) + "]"

	lines := []string{top}
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, textLine)
	if cfg.BlankRow {
		lines = append(lines, emptyInnerLine)
	}
	lines = append(lines, bottom)

	return lines, nil
}

