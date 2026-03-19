package main

import (
	"context"
	"fmt"
	"github.com/logicxking54/brain"
	"github.com/logicxking54/brain/utils"
	"log"
	"math/rand/v2"
	"os"
	"sync"
	"syscall"
	"time"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

type App struct {
	appCtx    brain.IContext
	ctx       context.Context
	ready     bool
	running   bool
	attacking bool
	f         *os.File
	logs      []map[string]string
	err       string
	once      sync.Once
	timeStr   string
}

func NewApp() *App {
	appCtx := brain.NewContext(&brain.ContextOptions{
		ENV: brain.NewEnv(),
	})

	return &App{
		appCtx: appCtx,
		logs:   make([]map[string]string, 0),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.ready = false
	a.attacking = false

	a.logs = append(a.logs, map[string]string{
		"date":    time.Now().Format(time.TimeOnly),
		"content": "กำลังบูสระบบ",
	})

	for i := 0; i < 3; i++ {
		f := findDevice(a)

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
		log.Println(a.err)
		return
	}

	a.logs = append(a.logs, map[string]string{
		"date":    time.Now().Format(time.TimeOnly),
		"content": "พบเครื่อง arduino แล้ว",
	})

	a.logs = append(a.logs, map[string]string{
		"date":    time.Now().Format(time.TimeOnly),
		"content": "ระบบพร้อมทำงาน",
	})
	a.ready = true
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

func (a *App) TimeStr() string {
	return a.timeStr
}

func (a *App) StartBot() {
	a.logs = append(a.logs, map[string]string{
		"date":    time.Now().Format(time.TimeOnly),
		"content": "เริ่มปล่อยบอท",
	})
	a.running = true
	a.attacking = true

	attack := func(a *App, ok *bool) {
		for utils.GetBool(ok) {
			pressDown(a, 't')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))
			pressUp(a, 't')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))

			pressDown(a, 'r')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))
			pressUp(a, 'r')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))

			pressDown(a, 'x')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))
			pressUp(a, 'x')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))
		}
	}

	collection := func(a *App, ok *bool) {
		for utils.GetBool(ok) {
			pressDown(a, 'f')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))
			pressUp(a, 'f')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))

			pressDown(a, 'w')
			time.Sleep(time.Second * time.Duration(rand.IntN(2-1+1)+1))
			pressUp(a, 'w')
			time.Sleep(time.Second)

			pressDown(a, 's')
			time.Sleep(time.Second * time.Duration(rand.IntN(2-1+1)+1))
			pressUp(a, 's')

			time.Sleep(time.Second * time.Duration(rand.IntN(2-1+1)+1))
		}
	}

	for a.running {
		go attack(a, &a.attacking)
		go collection(a, &a.attacking)

		start := time.Now()
		for a.running {
			if time.Since(start).Minutes() >= 10 && a.running {
				c := 0
				a.attacking = false

				for {
					time.Sleep(time.Second)
					c++

					if c >= 60 {
						break
					}

					if !a.running {
						break
					}
				}

				a.attacking = true
				break
			}

			a.timeStr = fmt.Sprintf("%dนาที %dวิ", int64(time.Since(start).Minutes()), int64(time.Since(start).Seconds()))
			time.Sleep(time.Second)
		}
	}
}

func (a *App) StartAutoKey() {
	a.logs = append(a.logs, map[string]string{
		"date":    time.Now().Format(time.TimeOnly),
		"content": "เปิดโหมด auto key",
	})
	a.running = true
	start := time.Now()

	for a.running {
		a.timeStr = fmt.Sprintf("%dนาที %dวิ", int64(time.Since(start).Minutes()), int64(time.Since(start).Seconds()))
		time.Sleep(time.Millisecond * 10)

		log.Println("getting keyboard key")
		ret, _, _ := getAsyncKeyState.Call(uintptr(0x54))

		if ret&0x8000 != 0 {

			pressDown(a, 'e')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(50-30+1)+30))
			pressUp(a, 'e')

			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (a *App) Stop() {
	a.logs = append(a.logs, map[string]string{
		"date":    time.Now().Format(time.TimeOnly),
		"content": "หยุดการทำงาน",
	})
	a.running = false
	a.attacking = false
}
