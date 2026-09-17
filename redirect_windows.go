//go:build windows

package main

import (
	"os"
	"syscall"
)

func isRedirect(path string, info os.FileInfo) (bool, error) {
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false, err
	}
	attributes, err := syscall.GetFileAttributes(p)
	if err != nil {
		return false, err
	}
	return attributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT != 0, nil
}
