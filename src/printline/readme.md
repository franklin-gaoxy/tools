# printline

It will automatically obtain the command window width and output a delimiter with the same width as the command window.

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

