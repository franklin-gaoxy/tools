//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

func showFatalError(err error) {
	if err == nil {
		return
	}
	messageBox("disk_not_sleep", err.Error())
}

func messageBox(title string, text string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("MessageBoxW")
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	textPtr, _ := syscall.UTF16PtrFromString(text)
	_, _, _ = proc.Call(0, uintptr(unsafe.Pointer(textPtr)), uintptr(unsafe.Pointer(titlePtr)), 0)
}
