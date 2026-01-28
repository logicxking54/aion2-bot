package main

import (
	"context"
	"fmt"
	"github.com/go-vgo/robotgo"
	"gitlab.logicxking.com/core/brain"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type App struct {
	appCtx      brain.IContext
	ctx         context.Context
	ready       bool
	running     bool
	f           *os.File
	logs        []map[string]string
	pid         int
	err         string
	overlayHwnd uintptr
	once        sync.Once
	monsterMats []MonsterMat
}

func NewApp() *App {
	appCtx := brain.NewContext(&brain.ContextOptions{
		ENV: brain.NewEnv(),
	})

	return &App{
		appCtx:      appCtx,
		logs:        make([]map[string]string, 0),
		monsterMats: make([]MonsterMat, 0),
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
	} else {
		a.logs = append(a.logs, map[string]string{
			"date":    time.Now().Format(time.TimeOnly),
			"content": "พบเครื่อง arduino แล้ว",
		})
	}

	for a.pid == 0 {
		a.logs = append(a.logs, map[string]string{
			"date":    time.Now().Format(time.TimeOnly),
			"content": "กำลังหาหน้าต่างเกม Aion2 ...",
		})

		windowNames, err := robotgo.FindNames()
		if err != nil {
			a.err = err.Error()
			a.Stop()
			log.Println(a.err)
			return
		}

		for i := range windowNames {
			if strings.ToLower(windowNames[i]) == "aion2.exe" {
				pids, err := robotgo.FindIds(windowNames[i])
				if err != nil {
					a.err = err.Error()
					a.Stop()
					log.Println(a.err)
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

		time.Sleep(time.Second * 2)
	}

	err := buildMonster(a)
	if err != nil {
		a.err = err.Error()
		a.Stop()
		log.Println(a.err)
		return
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
	a.running = true
	a.logs = append(a.logs, map[string]string{
		"date":    time.Now().Format(time.TimeOnly),
		"content": fmt.Sprintf("%v", a.running),
	})
}

func (a *App) Stop() {
	DrawOverlayRect(a, 0, 0, 0, 0)
	a.logs = append(a.logs, map[string]string{
		"date":    time.Now().Format(time.TimeOnly),
		"content": "หยุดการทำงาน",
	})
	a.ready = false
	a.running = false
	a.logs = append(a.logs, map[string]string{
		"date":    time.Now().Format(time.TimeOnly),
		"content": fmt.Sprintf("%v", a.running),
	})
}
