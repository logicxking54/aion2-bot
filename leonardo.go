package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
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

func clickDown(a *App) {
	a.f.Write([]byte{'D', 'L'})
}

func clickUp(a *App) {
	a.f.Write([]byte{'U', 'L'})
}

func moveMouse(a *App, targetX, targetY int) {
	steps := 800
	delay := time.Millisecond

	startX, startY := robotgo.Location()
	dx := float64(targetX-startX) / float64(steps)
	dy := float64(targetY-startY) / float64(steps)

	for i := 0; i < steps; i++ {
		x := int(float64(startX) + dx*float64(i))
		y := int(float64(startY) + dy*float64(i))
		robotgo.Move(x, y)
		time.Sleep(delay)
	}

	robotgo.Move(targetX, targetY)
}
