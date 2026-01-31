package main

import (
	"context"
	"fmt"
	"github.com/go-vgo/robotgo"
	"gitlab.logicxking.com/core/brain"
	"gitlab.logicxking.com/core/brain/utils"
	"log"
	"math/rand/v2"
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
	go func() {
		aa := 0
		for aa <= 100 {
			aa++
			time.Sleep(time.Second)
			x, y := robotgo.Location()
			log.Println(x, y, aa)
		}
	}()

	var displayRect robotgo.Rect
	startingUp := false
	respawning := false
	attacking := false

	startUp := func(a *App) {
		startingUp = true

		for {
			if !respawning && !attacking {
				break
			}

			log.Println("startUp checking", respawning, attacking)
			time.Sleep(time.Second)
		}

		displayRect = robotgo.GetDisplayRect(0)
		moveMouse(a, displayRect.X+(displayRect.W/2), displayRect.Y+(displayRect.H/2))
		log.Println("displayRect.X", displayRect.X)
		log.Println("displayRect.W", displayRect.W)
		log.Println("displayRect.Y", displayRect.Y)
		log.Println("displayRect.H", displayRect.H)
		clickDown(a)
		time.Sleep(time.Millisecond * 50)
		clickUp(a)

		for {
			time.Sleep(time.Second)
			log.Println("กำลังเปิด Map")
			ok, x1, y1, w, h := FindPosition(0, "./input/abyss2_menumap.png")
			if ok && w > 0 && h > 0 {
				log.Println("เปิด Map เรียบร้อย", ok, x1, y1, w, h)
				break
			} else {
				pressDown(a, 'm')
				time.Sleep(time.Millisecond * 50)
				pressUp(a, 'm')
			}
		}

		retry := 0
		moveMouse(a, displayRect.X+(displayRect.W/2), displayRect.H-100)

		for retry <= 4 {
			log.Println("กำลังหาจุดวาป")
			clickDown(a)
			moveMouse(a, displayRect.X+(displayRect.W/2), -(displayRect.H - 100))
			clickUp(a)

			ok, x, y, w, h := FindPosition(0, "./input/abyss2_teleport1.png")
			if ok && w > 0 && h > 0 {
				log.Println("เจอหาจุดวาปแล้ว")

				moveMouse(a, int(x), int(y))
				clickDown(a)
				time.Sleep(time.Millisecond * 50)
				clickUp(a)
				break
			}

			moveMouse(a, displayRect.X+(displayRect.W/2), displayRect.H-100)
			time.Sleep(time.Second)
			retry++
		}

		log.Println("กำลังกดปุ่มวาป")
		pressDown(a, 'f')
		time.Sleep(time.Millisecond * time.Duration(rand.IntN(500-100+1)+100))
		pressUp(a, 'f')
		time.Sleep(time.Millisecond * time.Duration(rand.IntN(500-100+1)+100))
		pressDown(a, 'f')
		time.Sleep(time.Millisecond * time.Duration(rand.IntN(500-100+1)+100))
		pressUp(a, 'f')
		time.Sleep(time.Millisecond * time.Duration(rand.IntN(1500-500+1)+500))

		moveMouse(a, displayRect.X+(displayRect.W/2), displayRect.Y+(displayRect.H/2))
		clickDown(a)
		time.Sleep(time.Millisecond * time.Duration(rand.IntN(500-100+1)+100))
		clickUp(a)
		time.Sleep(time.Millisecond * time.Duration(rand.IntN(1500-500+1)+500))

		log.Println("กำลังกดปุ่ม confirm")
		pressDown(a, 'f')
		time.Sleep(time.Millisecond * time.Duration(rand.IntN(500-100+1)+100))
		pressUp(a, 'f')
		pressDown(a, 'f')
		time.Sleep(time.Millisecond * time.Duration(rand.IntN(500-100+1)+100))
		pressUp(a, 'f')
		time.Sleep(time.Millisecond * time.Duration(rand.IntN(1500-500+1)+500))

		time.Sleep(time.Second * 5)
		startingUp = false
	}

	respawn := func(a *App, window *robotgo.Rect) {
		respawning = true

		for {
			if !startingUp && !attacking {
				break
			}

			log.Println("respawn checking", respawning, attacking)
			time.Sleep(time.Second)
		}

		moveMouse(a, window.X+(window.W/2), window.Y+(window.H/2)+75)
		clickDown(a)
		time.Sleep(time.Millisecond * 50)
		clickUp(a)

		respawning = false
	}

	attack := func(a *App, window *robotgo.Rect, ok *bool, k1 *bool, k2 *bool) {
		*ok = true

		for {
			if !utils.GetBool(k1) && !utils.GetBool(k2) {
				break
			}

			log.Println("attack checking", utils.GetBool(k1), utils.GetBool(k2))
			time.Sleep(time.Second)
		}

		for utils.GetBool(ok) {
			pressDown(a, 'r')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))
			pressUp(a, 'r')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))

			pressDown(a, 'e')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))
			pressUp(a, 'e')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))

			pressDown(a, '1')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))
			pressUp(a, '1')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))

			pressDown(a, '2')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))
			pressUp(a, '2')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))

			time.Sleep(time.Second * time.Duration(rand.IntN(2-1+1)+1))
		}
	}

	collection := func(a *App, ok *bool, k1 *bool, k2 *bool) {
		*ok = true

		for {
			if !utils.GetBool(k1) && !utils.GetBool(k2) {
				break
			}

			log.Println("attack checking", utils.GetBool(k1), utils.GetBool(k2))
			time.Sleep(time.Second)
		}

		for utils.GetBool(ok) {
			pressDown(a, 'f')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))
			pressUp(a, 'f')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+30))

			pressDown(a, 'w')
			time.Sleep(time.Second * time.Duration(rand.IntN(4-1+1)+1))
			pressUp(a, 'w')
			time.Sleep(time.Second)

			pressDown(a, 's')
			time.Sleep(time.Second * time.Duration(rand.IntN(4-1+1)+1))
			pressUp(a, 's')

			time.Sleep(time.Second * time.Duration(rand.IntN(4-2+1)+2))
		}
	}

	isDead := func() bool {
		die, _, _, _, _ := FindPosition(0, "./input/dead.png")
		if die {
			log.Println("คุณตายแล้ว")
			return true
		} else {
			log.Println("คุณยังไม่ตาย")
			return false
		}
	}

	rootDeadCount := 0
	for {
		startUp(a)
		go attack(a, &displayRect, &attacking, &startingUp, &respawning)
		go collection(a, &attacking, &startingUp, &respawning)

		for {
			if isDead() {
				rootDeadCount++

				if rootDeadCount >= 5 {
					attacking = false
					respawn(a, &displayRect)
					rootDeadCount = 0
					time.Sleep(time.Minute * 3)
					break
				}
			}

			time.Sleep(time.Second)
		}
	}
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
