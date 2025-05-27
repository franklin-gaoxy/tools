# printline

It will automatically obtain the command window width and output a delimiter with the same width as the command window.

## golang edition

> This is the first recommended.

Uses the cobra framework, supports line and center modes, and also supports other parameters. Use -h to view details.

### demo

```shell
printline line -
printline line -l 3 - # Output three lines
printline center "hello world" # Print "hello world" in the center
./printline center -b y "hello world" -s - -p y # Print in the center, specify symbols, enable detailed information output, print blank lines at the beginning and end
```

### build

```shell
go build -o printline main.go
```





## python edition

code file: [printline.py](https://github.com/franklin-gaoxy/tools/blob/main/src/printline/printline.py)

### instructions

parameters: center ,Center display of specified characters.

parameters: line, Output the specified line of characters.

### demo

```shell
printline # print help
printline line = # print one line =
printline line "   " # print there blank line
printline center 'hello world' # print centered 'hello world',default line symbol is =
printline center 'hello world' - # print centered 'hello world',line sumbol is -
```

### build

```
pip install pyinstaller
pyinstaller --onefile printline.py
```

