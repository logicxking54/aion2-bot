package main

import (
	"context"
	"fmt"
	"github.com/go-vgo/robotgo"
	"os"
	"strings"
	"time"
)

type App struct {
	ctx     context.Context
	ready   bool
	running bool
	f       *os.File
	logs    []map[string]string
	pid     int
	err     string
}

func NewApp() *App {
	return &App{
		logs: make([]map[string]string, 0),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.ready = false

	a.logs = append(a.logs, map[string]string{
		"date":    time.Now().Format(time.TimeOnly),
		"content": "กำลังบูสระบบ",
	})

	for i := 0; i < 3; i++ {
		f := a.findDevice()

		if f == nil {
			a.logs = append(a.logs, map[string]string{
				"date":    time.Now().Format(time.TimeOnly),
				"content": "ไม่พบเครื่อง arduino, กำลังรอใหม่อีกครั้ง ...",
			})
			continue
		}

		a.f = f
		break
	}

	if a.f == nil {
		a.logs = append(a.logs, map[string]string{
			"date":    time.Now().Format(time.TimeOnly),
			"content": "ไม่พบเครื่อง arduino, กรุณาเชื่อมต่อกับเครื่องก่อน แล้วลองอีกครั้ง",
		})
		a.err = "ไม่พบเครื่อง arduino, กรุณาเชื่อมต่อกับเครื่องก่อน แล้วลองอีกครั้ง"
		return
	} else {
		a.logs = append(a.logs, map[string]string{
			"date":    time.Now().Format(time.TimeOnly),
			"content": "พบเครื่อง arduino แล้ว",
		})
	}

	a.logs = append(a.logs, map[string]string{
		"date":    time.Now().Format(time.TimeOnly),
		"content": "ระบบพร้อมทำงาน",
	})
	a.ready = true
	go startLoop(a)
}

func (a *App) IsRunning() bool {
	return a.running
}

func (a *App) Log() []map[string]string {
	return a.logs
}

func (a *App) Error() string {
	return a.err
}

func (a *App) Start() {
	a.kill()
	a.running = true
}

func (a *App) Stop() {
	a.kill()
	a.logs = append(a.logs, map[string]string{
		"date":    time.Now().Format(time.TimeOnly),
		"content": "หยุดการทำงาน",
	})
	a.ready = false
	a.running = false
}

func startLoop(a *App) {
	for {
		time.Sleep(time.Second)

		if a.running {
			for a.pid == 0 {
				a.logs = append(a.logs, map[string]string{
					"date":    time.Now().Format(time.TimeOnly),
					"content": "กำลังหาหน้าต่างเกม Aion2 ...",
				})

				windowNames, err := robotgo.FindNames()
				if err != nil {
					a.err = err.Error()
					a.Stop()
					return
				}

				for i := range windowNames {
					if strings.ToLower(windowNames[i]) == "aion2.exe" {
						pids, err := robotgo.FindIds(windowNames[i])
						if err != nil {
							a.err = err.Error()
							a.Stop()
							return
						}

						a.pid = pids[0]
						a.logs = append(a.logs, map[string]string{
							"date":    time.Now().Format(time.TimeOnly),
							"content": fmt.Sprintf("พบหน้าต่าง Aion2 แล้ว (%d)", a.pid),
						})
						break
					}
				}

				a.logs = append(a.logs, map[string]string{
					"date":    time.Now().Format(time.TimeOnly),
					"content": "ไม่พบหน้าต่างเกม Aion2, กำลังลองใหม่อีกครั้ง ...",
				})
				time.Sleep(time.Second * 2)
			}

			c := 0
			for a.running {
				if a.pid != robotgo.GetPid() {
					a.logs = append(a.logs, map[string]string{
						"date":    time.Now().Format(time.TimeOnly),
						"content": "กรุณาเปิดหน้าต่าง Aion2, โปรแกรมไม่สามารถทำงานได้ หากไม่โฟกัสที่หน้าต่าง Aion2 ตลอดเวลา",
					})
					time.Sleep(time.Second)
					continue
				}

				a.logs = append(a.logs, map[string]string{
					"date":    time.Now().Format(time.TimeOnly),
					"content": "กด W",
				})
				a.pressDown("w")

				time.Sleep(time.Second)

				a.pressUp("w")
				a.logs = append(a.logs, map[string]string{
					"date":    time.Now().Format(time.TimeOnly),
					"content": "ปล่อย W",
				})

				if c >= 3 {
					break
				}
				c++
			}

			a.kill()
		}
	}
}

func (a *App) findDevice() *os.File {
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

func (a *App) pressDown(k string) {
	a.f.Write([]byte(fmt.Sprintf("D|%s", k) + "\n"))
}

func (a *App) pressUp(k string) {
	a.f.Write([]byte(fmt.Sprintf("U|%s", k) + "\n"))
}

func (a *App) kill() {
	a.f.Write([]byte("KILL"))
}
