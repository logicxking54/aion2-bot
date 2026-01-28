package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"gocv.io/x/gocv"
	"log"
	"time"
)

func startLoop(a *App) {
	for {
		time.Sleep(time.Second)

		if a.running {
			err := robotgo.ActivePid(a.pid)
			if err != nil {
				a.err = err.Error()
				a.Stop()
				log.Println(a.err)
				return
			}

			keyboard, langId := GetCurrentKeyboardLayout(a)
			a.logs = append(a.logs, map[string]string{
				"date":    time.Now().Format(time.TimeOnly),
				"content": fmt.Sprintf("คีบอร์ทกำลังใช้ %s (%d)", keyboard, langId),
			})

			if robotgo.GetPid() != a.pid {
				a.logs = append(a.logs, map[string]string{
					"date":    time.Now().Format(time.TimeOnly),
					"content": "กรุณาเปิดหน้าต่าง Aion2, โปรแกรมไม่สามารถทำงานได้ หากไม่โฟกัสที่หน้าต่าง Aion2 ตลอดเวลา",
				})
				time.Sleep(time.Second)
				continue
			}

			for a.running {
				if a.pid != robotgo.GetPid() {
					a.logs = append(a.logs, map[string]string{
						"date":    time.Now().Format(time.TimeOnly),
						"content": "กรุณาเปิดหน้าต่าง Aion2, โปรแกรมไม่สามารถทำงานได้ หากไม่โฟกัสที่หน้าต่าง Aion2 ตลอดเวลา",
					})
					time.Sleep(time.Second)
					continue
				}

				target := gocv.NewMat()
				worldBit := robotgo.CaptureScreen()
				worldImage := robotgo.ToImage(worldBit)
				robotgo.FreeBitmap(worldBit)

				worldMat, err := gocv.ImageToMatRGB(worldImage)
				if err != nil {
					a.err = err.Error()
					a.Stop()
					log.Println(a.err)
					return
				}

				for i := range a.monsterMats {
					err = gocv.MatchTemplate(worldMat, a.monsterMats[i].Mat, &target, gocv.TmCcoeffNormed, gocv.NewMat())
					if err != nil {
						a.err = err.Error()
						a.Stop()
						return
					}

					_, maxVal, _, maxLoc := gocv.MinMaxLoc(target)
					log.Println("cond", maxVal)
					if maxVal >= 0.75 {
						DrawOverlayRect(a, maxLoc.X, maxLoc.Y, a.monsterMats[i].Mat.Cols(), a.monsterMats[i].Mat.Rows())
					}
				}

				time.Sleep(time.Millisecond * 100)
			}

			a.Stop()
		}
	}
}
