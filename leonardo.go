package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"math"
	"math/rand"
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
	clamp := func(v, min, max int) int {
		if v < min {
			return min
		}
		if v > max {
			return max
		}
		return v
	}

	curX, curY := robotgo.Location()
	dx := targetX - curX
	dy := targetY - curY

	steps := int(math.Hypot(float64(dx), float64(dy))/6) + 12
	if steps < 20 {
		steps = 20
	}

	var lastX, lastY float64

	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		ease := t * t * (3 - 2*t)

		curXf := float64(dx) * ease
		curYf := float64(dy) * ease

		moveX := int(curXf - lastX)
		moveY := int(curYf - lastY)

		lastX = curXf
		lastY = curYf

		moveX = clamp(moveX, -127, 127)
		moveY = clamp(moveY, -127, 127)

		a.f.Write([]byte{'M', byte(int8(moveX)), byte(int8(moveY))})
		time.Sleep(time.Duration(rand.Intn(5)+4) * time.Millisecond)
	}
}
