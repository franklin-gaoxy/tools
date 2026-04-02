//go:build !windows

package main

import (
	"fmt"
	"os"
)

func showFatalError(err error) {
	if err == nil {
		return
	}
	_, _ = fmt.Fprintln(os.Stderr, err.Error())
}
