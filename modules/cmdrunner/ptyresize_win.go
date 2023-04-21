//go:build windows

package cmdrunner

import (
	"os"
)

func resizePty(f *os.File) func() {
	return func() {}
}
