//go:build !windows

package cmdrunner

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/wtfutil/wtf/logger"
)

func resizePty(f *os.File) func() {
	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, f); err != nil {
				logger.Log(fmt.Sprintf("error resizing pty: %s", err))
			}
		}
	}()
	ch <- syscall.SIGWINCH                       // Initial resize.
	return func() { signal.Stop(ch); close(ch) } // Cleanup signals when done.
}
