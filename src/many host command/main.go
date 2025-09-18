package main

import "many/tools"

func main() {
	rootcmd := tools.BoundBobraAgrs()
	tools.ExecuteStart(rootcmd)
}
