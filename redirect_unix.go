//go:build !windows

package main

import "os"

func isRedirect(path string, info os.FileInfo) (bool, error) {
	return info.Mode()&os.ModeSymlink != 0, nil
}
