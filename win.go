package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	gdi32                   = syscall.NewLazyDLL("gdi32.dll")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadPID  = user32.NewProc("GetWindowThreadProcessId")
	procGetKeyboardLayout   = user32.NewProc("GetKeyboardLayout")
)

func GetCurrentKeyboardLayout(a *App) (string, uint16) {
	hwnd, _, _ := procGetForegroundWindow.Call()

	var pid uint32
	tid, _, _ := procGetWindowThreadPID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))

	hkl, _, _ := procGetKeyboardLayout.Call(tid)
	langID := uint16(hkl & 0xFFFF)

	switch langID {
	case 0x0409:
		return "English (US)", langID
	case 0x041E:
		return "Thai", langID
	default:
		return fmt.Sprintf("Unknown (0x%04X)", langID), langID
	}
}

func DrawOverlayRect(a *App, x, y, w, h int) {
	a.once.Do(func() {
		createWindowEx := user32.NewProc("CreateWindowExW")
		showWindow := user32.NewProc("ShowWindow")
		setLayered := user32.NewProc("SetLayeredWindowAttributes")
		getSystemMetrics := user32.NewProc("GetSystemMetrics")
		setWindowPos := user32.NewProc("SetWindowPos")

		const (
			WS_EX_LAYERED     = 0x80000
			WS_EX_TRANSPARENT = 0x20
			WS_EX_TOPMOST     = 0x8
			WS_EX_NOACTIVATE  = 0x08000000
			WS_EX_TOOLWINDOW  = 0x00000080

			WS_POPUP     = 0x80000000
			SW_SHOW      = 5
			LWA_COLORKEY = 0x1

			SM_CXSCREEN = 0
			SM_CYSCREEN = 1

			HWND_TOPMOST   = ^uintptr(0)
			SWP_NOMOVE     = 0x2
			SWP_NOSIZE     = 0x1
			SWP_NOACTIVATE = 0x10
		)

		className, _ := syscall.UTF16PtrFromString("STATIC")

		wScr, _, _ := getSystemMetrics.Call(SM_CXSCREEN)
		hScr, _, _ := getSystemMetrics.Call(SM_CYSCREEN)

		hwnd, _, _ := createWindowEx.Call(
			WS_EX_LAYERED|WS_EX_TRANSPARENT|WS_EX_TOPMOST|WS_EX_NOACTIVATE|WS_EX_TOOLWINDOW,
			uintptr(unsafe.Pointer(className)),
			uintptr(unsafe.Pointer(className)),
			WS_POPUP,
			0, 0,
			wScr, hScr,
			0, 0, 0, 0,
		)

		transparentColor := uint32(0x000000)
		setLayered.Call(hwnd, uintptr(transparentColor), 0, LWA_COLORKEY)

		setWindowPos.Call(
			hwnd,
			HWND_TOPMOST,
			0, 0, 0, 0,
			SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE,
		)

		showWindow.Call(hwnd, SW_SHOW)
		a.overlayHwnd = hwnd
	})

	getDC := user32.NewProc("GetDC")
	releaseDC := user32.NewProc("ReleaseDC")

	createPen := gdi32.NewProc("CreatePen")
	selectObj := gdi32.NewProc("SelectObject")
	deleteObj := gdi32.NewProc("DeleteObject")
	rectangle := gdi32.NewProc("Rectangle")
	getStockObject := gdi32.NewProc("GetStockObject")
	patBlt := gdi32.NewProc("PatBlt")
	createSolidBrush := gdi32.NewProc("CreateSolidBrush")

	const (
		PS_SOLID   = 0
		NULL_BRUSH = 5
		BLACKNESS  = 0x00000042
	)

	hdc, _, _ := getDC.Call(a.overlayHwnd)
	if hdc == 0 {
		return
	}
	defer releaseDC.Call(a.overlayHwnd, hdc)

	// ===== CLEAR OLD DRAW =====
	brushClear, _, _ := createSolidBrush.Call(0x000000)
	oldBrush2, _, _ := selectObj.Call(hdc, brushClear)

	patBlt.Call(hdc, 0, 0, 5000, 5000, BLACKNESS)

	selectObj.Call(hdc, oldBrush2)
	deleteObj.Call(brushClear)

	// ===== DRAW NEW RECT =====
	red := uint32(0x0000FF)
	pen, _, _ := createPen.Call(PS_SOLID, 3, uintptr(red))
	oldPen, _, _ := selectObj.Call(hdc, pen)

	brush, _, _ := getStockObject.Call(NULL_BRUSH)
	oldBrush, _, _ := selectObj.Call(hdc, brush)

	rectangle.Call(
		hdc,
		uintptr(x),
		uintptr(y),
		uintptr(x+w),
		uintptr(y+h),
	)

	selectObj.Call(hdc, oldBrush)
	selectObj.Call(hdc, oldPen)
	deleteObj.Call(pen)
}
