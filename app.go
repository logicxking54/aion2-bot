package main

import (
	"context"
	"os"
	"time"
)

type App struct {
	ctx     context.Context
	running bool
	logs    []map[string]string
}

func NewApp() *App {
	return &App{
		logs: []map[string]string{
			{
				"date":    time.Now().Format(time.TimeOnly),
				"content": "ระบบพร้อมทำงาน",
			},
		},
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.startLoop()
}

func (a *App) IsRunning() bool {
	return a.running
}

func (a *App) Log() []map[string]string {
	return a.logs
}

func (a *App) Start() {
	a.running = true
}

func (a *App) Stop() {
	a.running = false
}

func (a *App) startLoop() {
	for {
		if a.running {
			a.logs = append(a.logs, map[string]string{
				"date":    time.Now().Format(time.TimeOnly),
				"content": "กำลังเริ่มทำงาน",
			})

			f, err := os.OpenFile("COM6", os.O_RDWR, 0)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			time.Sleep(2 * time.Second)

			for i := 0; i < 10; i++ {
				f.Write([]byte("W"))

				a.logs = append(a.logs, map[string]string{
					"date":    time.Now().Format(time.TimeOnly),
					"content": "กด W",
				})

				time.Sleep(time.Second)
			}
		} else {
			time.Sleep(time.Second)
		}
	}
}
