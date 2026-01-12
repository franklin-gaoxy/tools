# Printline

A CLI tool to generate and print separator lines or centered text in the terminal. Useful for scripts and logs formatting.

## Build

```bash
# Using Make
make build

# Using Go directly
go build -o printline .

# Or build only the entry file (also supported)
go build -o printline main.go
```

## Usage

```bash
./printline [command] [flags]
```

---

### Command: `line`

Print a line of specified characters.

#### Function
Repeats a character (or string) across the entire width of the terminal.

#### Flags

##### `-l, --lines int`
Number of lines.
- **Default**: 1
- **Function**: Specifies how many times to repeat the line vertically.

##### `-b, --blank-row`
Print blank rows.
- **Type**: bool flag
- **Default**: disabled
- **Function**: If present, adds an empty line before and after the printed line for spacing.

#### Examples
```bash
# Print a line of '='
./printline line

# Print a line of '*'
./printline line "*"

# Print 2 lines with blank rows
./printline line -l 2 -b
```

---

### Command: `center`

Print text centered in the terminal.

#### Function
Calculates the terminal width and prints the provided text in the middle, padding the left and right with a specified symbol.

#### Flags

##### `-s, --symbol string`
Padding symbol.
- **Default**: `=`
- **Function**: The character used to fill the empty space around the centered text.

##### `-p, --print-info`
Print debug info.
- **Type**: bool flag
- **Function**: If present, prints details about terminal width and string length calculations.

##### `-b, --blank-row`
Print blank rows.
- **Type**: bool flag
- **Default**: disabled
- **Function**: Adds vertical padding (empty lines) around the output.

#### Examples
```bash
# Print "HELLO" centered with '='
./printline center "HELLO"

# Print "WARNING" centered with '!'
./printline center "WARNING" -s "!" -b -p
```

---

### Command: `version`

Print version information.

#### Function
Displays the current version of the tool.

---

### Command: `completely-center`

Print text inside a full-width border.

#### Function
Draws a border whose width matches the terminal width, and prints the provided text inside the border.

#### Args
Requires exactly 1 argument: the text to print.

#### Flags

##### `--style string`
Border style.
- **Default**: `box`
- **Values**: `box`, `ascii`, `solid`, `double`
- **Function**: Selects the base border rendering style.

##### `-t, --template string`
Template key.
- **Default**: empty (uses `--style`)
- **Function**: Overrides `--style` and selects a template by key. Built-in keys include: `box`, `ascii`, `solid`, `double`, `log`, `warn`, `error`, `success`. Intended for extending with more templates in code.

##### `-s, --symbol string`
Border symbol (solid style).
- **Default**: `=`
- **Function**: Used only when `--style solid` (or when `-s` is provided without specifying `--style`).

#### Notes
- If you provide `-s/--symbol` without `--style`, it will automatically switch to `solid`.
- If you provide `-t/--template`, it takes precedence over `--style`.

##### `-h, --header string`
Line prefix.
- **Default**: empty
- **Function**: Printed before every output line. The border width will be reduced to keep total width within the terminal.

##### `-b, --blank-row`
Add blank rows.
- **Type**: bool flag
- **Default**: disabled
- **Function**: If present, adds an empty line above and below the text inside the border.

#### Examples
```bash
# Basic
./printline completely-center "Hello World"

# ASCII style (|, -, +)
./printline completely-center "Hello World" --style ascii

# Double-line border template
./printline completely-center "Hello World" -t double

# Log / Warn / Error / Success templates
./printline completely-center "Something happened" -t log
./printline completely-center "Disk usage high" -t warn
./printline completely-center "Failed to connect" -t error
./printline completely-center "Deployment finished" -t success

# Use '!' as border symbol
./printline completely-center "Hello World" --style solid -s "!"

# Add a prefix for each line
./printline completely-center "Hello World" -h "# "

# Add blank rows inside the border
./printline completely-center "Hello World" -b
```
