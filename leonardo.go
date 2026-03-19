package main

import (
	"fmt"
	"os"
	"time"
)

func findDevice(a *App) *os.File {
	for i := 1; i <= 256; i++ {
		port := fmt.Sprintf("COM%d", i)

		f, err := os.OpenFile(port, os.O_RDWR, 0)
		if err == nil {
			time.Sleep(2 * time.Second)
			return f
		}
	}

	return nil
}

func pressDown(a *App, k byte) {
	a.f.Write([]byte{'D', k})
}

func pressUp(a *App, k byte) {
	a.f.Write([]byte{'U', k})
}
