package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/google/gopacket/pcap"
	"github.com/logicxking54/brain"
	"github.com/logicxking54/brain/utils"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

type App struct {
	appCtx      brain.IContext
	ctx         context.Context
	ready       bool
	running     bool
	attacking   bool
	f           *os.File
	logs        []map[string]string
	err         string
	once        sync.Once
	pcapHanlder *pcap.Handle
	timeStr     string
	eUseAt      *time.Time
	rUseAt      *time.Time
	r3UseAt     *time.Time
	tUseAt      *time.Time
	eCooldown   int64
	lCooldown   int64
	rCooldown   int64
}

func NewApp() *App {
	appCtx := brain.NewContext(&brain.ContextOptions{
		ENV: brain.NewEnv(),
	})

	return &App{
		appCtx:    appCtx,
		logs:      make([]map[string]string, 0),
		eUseAt:    utils.ToPointer(time.Now()),
		eCooldown: 3000,
		lCooldown: 800,
		rCooldown: 800,
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

func (a *App) ECooldown() string {
	if a.eUseAt == nil {
		return "พร้อมใช้งาน"
	}

	if (a.eCooldown-time.Since(*a.eUseAt).Milliseconds())/1000 < 0 {
		return "พร้อมใช้งาน"
	}

	return fmt.Sprintf("%d วินาที", (a.eCooldown-time.Since(*a.eUseAt).Milliseconds())/1000)
}

func (a *App) LCooldown() string {
	return "พร้อมใช้งาน"
}

func (a *App) RCooldown() string {
	return "พร้อมใช้งาน"
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

		ret, _, _ := getAsyncKeyState.Call(uintptr(0xDB))
		if ret&0x8000 != 0 {
			pressDown(a, 'e')
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(50-30+1)+30))
			pressUp(a, 'e')

			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (a *App) StartAutoKeyNetwork() {
	//devices, err := pcap.FindAllDevs()
	//if err != nil || len(devices) <= 0 {
	//	a.logs = append(a.logs, map[string]string{
	//		"date":    time.Now().Format(time.TimeOnly),
	//		"content": "ค้นหา network device ล้มเหลว",
	//	})
	//	a.err = "ค้นหา network device ล้มเหลว"
	//	log.Println(a.err)
	//	return
	//}
	//
	//var deviceName string
	//for _, device := range devices {
	//	if strings.Contains(strings.ToLower(device.Description), "loopback") {
	//		deviceName = device.Name
	//		a.logs = append(a.logs, map[string]string{
	//			"date":    time.Now().Format(time.TimeOnly),
	//			"content": fmt.Sprintf("network device = %s", deviceName),
	//		})
	//		break
	//	}
	//}
	//
	//ih, err := pcap.NewInactiveHandle(deviceName)
	//if err != nil {
	//	a.logs = append(a.logs, map[string]string{
	//		"date":    time.Now().Format(time.TimeOnly),
	//		"content": err.Error(),
	//	})
	//	a.err = err.Error()
	//	log.Println(a.err)
	//	return
	//}
	//
	//ih.SetSnapLen(256)
	//ih.SetPromisc(true)
	//ih.SetTimeout(time.Millisecond)
	//ih.SetImmediateMode(true)
	//
	//handle, err := ih.Activate()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = handle.SetBPFFilter("tcp and host 127.0.0.1 and port 52744")
	//if err != nil {
	//	a.logs = append(a.logs, map[string]string{
	//		"date":    time.Now().Format(time.TimeOnly),
	//		"content": err.Error(),
	//	})
	//	a.err = err.Error()
	//	log.Println(a.err)
	//	return
	//}
	//
	a.logs = append(a.logs, map[string]string{
		"date":    time.Now().Format(time.TimeOnly),
		"content": "เปิดโหมด auto key",
	})
	//a.pcapHanlder = handle
	a.running = true
	//
	//go func(running *bool) {
	//	for *running {
	//		data, _, err := a.pcapHanlder.ZeroCopyReadPacketData()
	//		if err != nil {
	//			continue
	//		}
	//
	//		a.processGamePacket(data)
	//	}
	//}(&a.running)
	//
	//go func(running *bool) {
	//	for *running {
	//		ret, _, _ := getAsyncKeyState.Call(uintptr(0xDB))
	//		if ret&0x8000 != 0 {
	//			if a.rUseAt == nil {
	//				pressDown(a, 'r')
	//				time.Sleep(time.Millisecond * time.Duration(rand.IntN(10-3+1)+3))
	//				pressUp(a, 'r')
	//			} else if a.tUseAt == nil {
	//				pressDown(a, 't')
	//				time.Sleep(time.Millisecond * time.Duration(rand.IntN(10-3+1)+3))
	//				pressUp(a, 't')
	//			}
	//		}
	//
	//		time.Sleep(time.Millisecond)
	//	}
	//}(&a.running)
	//
	go func(running *bool) {
		for *running {
			ret, _, _ := getAsyncKeyState.Call(uintptr(0xDB))
			if ret&0x8000 != 0 {
				pressDown(a, 'r')
				time.Sleep(time.Millisecond * time.Duration(rand.IntN(10-3+1)+3))
				pressUp(a, 'r')
				time.Sleep(time.Millisecond * time.Duration(rand.IntN(80-30+1)+70))
			}

			time.Sleep(time.Millisecond)
		}
	}(&a.running)

	go func(running *bool) {
		for *running {
			ret, _, _ := getAsyncKeyState.Call(uintptr(0xDB))
			if ret&0x8000 != 0 {
				pressDown(a, 'e')
				time.Sleep(time.Millisecond * time.Duration(rand.IntN(50-10+1)+10))
				pressUp(a, 'e')
				time.Sleep(time.Millisecond * time.Duration(rand.IntN(120-30+1)+30))
			}

			time.Sleep(time.Millisecond)
		}
	}(&a.running)
}

func (a *App) readVarInt(data []byte, offset *int) int {
	value := 0
	shift := uint(0)
	for {
		if *offset >= len(data) {
			return -1
		}
		b := int(data[*offset])
		*offset++
		value |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			break
		}
		shift += 7
		if shift >= 32 {
			return -1
		}
	}
	return value
}

func (a *App) processGamePacket(data []byte) {
	if len(data) < 15 {
		return
	}

	idx := bytes.Index(data, []byte{0x04, 0x38})
	if idx == -1 {
		return
	}

	curr := idx + 2
	_ = a.readVarInt(data, &curr)
	_ = a.readVarInt(data, &curr)
	_ = a.readVarInt(data, &curr)
	_ = a.readVarInt(data, &curr)

	if curr+4 <= len(data) {
		skillRaw := binary.LittleEndian.Uint32(data[curr : curr+4])
		skillID := int(skillRaw) - (int(skillRaw) % 10000)

		switch skillID {
		case 1010000:
			return

		case 15050000: // Blaze
			if skillRaw == 15053450 {
				a.eUseAt = utils.ToPointer(time.Now())
			}

		case 15210000: // Flame Arrow
			a.rUseAt = utils.ToPointer(time.Now())
			a.tUseAt = nil
			log.Println("Flame Arrow", skillID, skillRaw, time.Now())

		case 15030000: // Burst
			a.rUseAt = utils.ToPointer(time.Now())
			a.tUseAt = nil
			log.Println("Burst", skillID, skillRaw, time.Now())

		case 15250000: // Pyroclasm
			if skillRaw == 15250450 {
				now := time.Now()
				if a.r3UseAt != nil && now.Sub(*a.r3UseAt) < time.Millisecond {
					return
				}
				a.rUseAt = utils.ToPointer(now)
				a.r3UseAt = utils.ToPointer(now)
				a.tUseAt = nil
				a.eUseAt = nil
				log.Println("Pyroclasm", skillID, skillRaw, now)
			}

		case 15090000: // Ice Chain
			if skillRaw == 15092340 {
				now := time.Now()
				if a.tUseAt != nil && now.Sub(*a.tUseAt) < time.Millisecond {
					return
				}
				a.tUseAt = utils.ToPointer(now)
				a.rUseAt = nil
				log.Println("Ice Chain", skillID, skillRaw, now)
			}

		case 15100000: // Cold Wave
			now := time.Now()
			if a.tUseAt != nil && now.Sub(*a.tUseAt) < time.Millisecond {
				return
			}
			a.tUseAt = utils.ToPointer(now)
			a.rUseAt = nil
			log.Println("Cold Wave", skillID, skillRaw, now)
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

	if a.pcapHanlder != nil {
		a.pcapHanlder.Close()
	}

	pressUp(a, 'r')
	pressUp(a, 't')
	pressUp(a, 'e')
}
