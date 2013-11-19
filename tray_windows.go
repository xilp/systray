package systray

import (
	"errors"
	"path/filepath"
	"unsafe"
	"syscall"
)

func (p *_Systray) Stop() error {
	nid := NOTIFYICONDATA{
		UID:  p.id,
		HWnd: HWND(p.hwnd),
	}
	nid.CbSize = uint32(unsafe.Sizeof(nid))

	ret, _, _ := Shell_NotifyIcon.Call(NIM_DELETE, uintptr(unsafe.Pointer(&nid)))
	if ret == 0 {
		return errors.New("shell notify delete failed")
	}
	return nil
}

func (p *_Systray) Show(file string, hint string) error {
	path := filepath.Join(p.iconPath, file)
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	icon, err := NewIconFromFile(path)
	if err != nil {
		return err
	}
	err = p.SetIcon(HICON(icon))
	if err != nil {
		return err
	}
	err = p.SetTooltip(hint)
	if err != nil {
		return err
	}
	return p.SetVisible(true)
}

func (p *_Systray) OnClick(fun func()) {
	p.lclick = fun
	p.rclick = fun
	p.dclick = fun
}

func (p *_Systray) SetTooltip(tooltip string) error {
	nid := NOTIFYICONDATA{
		UID:  p.id,
		HWnd: HWND(p.hwnd),
	}
	nid.CbSize = uint32(unsafe.Sizeof(nid))

	nid.UFlags = NIF_TIP
	copy(nid.SzTip[:], syscall.StringToUTF16(tooltip))

	ret, _, _ := Shell_NotifyIcon.Call(NIM_MODIFY, uintptr(unsafe.Pointer(&nid)))
	if ret == 0 {
		return errors.New("shell notify tooltip failed")
	}
	return nil
}

func (p *_Systray) SetVisible(visible bool) error {
	nid := NOTIFYICONDATA{
		UID:  p.id,
		HWnd: HWND(p.hwnd),
	}
	nid.CbSize = uint32(unsafe.Sizeof(nid))

	nid.UFlags = NIF_STATE
	nid.DwStateMask = NIS_HIDDEN
	if !visible {
		nid.DwState = NIS_HIDDEN
	}

	ret, _, _ := Shell_NotifyIcon.Call(NIM_MODIFY, uintptr(unsafe.Pointer(&nid)))
	if ret == 0 {
		return errors.New("shell notify tooltip failed")
	}
	return nil
}

func (p *_Systray) SetIcon(hicon HICON) error {
	nid := NOTIFYICONDATA{
		UID:  p.id,
		HWnd: HWND(p.hwnd),
	}
	nid.CbSize = uint32(unsafe.Sizeof(nid))

	nid.UFlags = NIF_ICON
	if hicon == 0 {
		nid.HIcon = 0
	} else {
		nid.HIcon = hicon
	}

	ret, _, _ := Shell_NotifyIcon.Call(NIM_MODIFY, uintptr(unsafe.Pointer(&nid)))
	if ret == 0 {
		return errors.New("shell notify icon failed")
	}
	return nil
}

func (p *_Systray) WinProc(hwnd HWND, msg uint32, wparam, lparam uintptr) uintptr {
	if msg == NotifyIconMessageId {
		if lparam == WM_LBUTTONDBLCLK {
			p.dclick()
		} else if (lparam == WM_LBUTTONUP) {
			p.lclick()
		} else if (lparam ==  WM_RBUTTONUP) {
			p.rclick()
		}
	}
	result, _, _ := DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wparam, lparam)
	return result
}

func (p *_Systray) Run() error {
	hwnd := p.mhwnd
	var msg MSG
	for {
		rt, _, _ := GetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		switch int(rt) {
		case 0:
			return nil
		case -1:
			return errors.New("run failed")
		}

		is, _, _ := IsDialogMessage.Call(hwnd, uintptr(unsafe.Pointer(&msg)))
		if is == 0 {
			TranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			DispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
		}
	}
	return nil
}

func _NewSystray(iconPath string, clientPath string, port int) *_Systray {
	tray, err := _NewSystrayEx(iconPath)
	if err != nil {
		panic(err)
	}
	return tray
}

func _NewSystrayEx(iconPath string) (*_Systray, error) {
	ni := &_Systray{iconPath, 0, 0, 0, func(){}, func(){}, func(){}}

	MainClassName := "MainForm"
	RegisterWindow(MainClassName, ni.WinProc)

	mhwnd, _, _ := CreateWindowEx.Call(
		WS_EX_CONTROLPARENT,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(MainClassName))),
		0,
		WS_OVERLAPPEDWINDOW|WS_CLIPSIBLINGS,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		0,
		0,
		0,
		0)
	if mhwnd == 0 {
		return nil, errors.New("create main win failed")
	}

	NotifyIconClassName := "NotifyIconForm"
	RegisterWindow(NotifyIconClassName, ni.WinProc)

	hwnd, _, _ := CreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(NotifyIconClassName))),
		0,
		0,
		0,
		0,
		0,
		0,
		uintptr(HWND_MESSAGE),
		0,
		0,
		0)
	if hwnd == 0 {
		return nil, errors.New("create notify win failed")
	}

	nid := NOTIFYICONDATA{
		HWnd:             HWND(hwnd),
		UFlags:           NIF_MESSAGE | NIF_STATE,
		DwState:          NIS_HIDDEN,
		DwStateMask:      NIS_HIDDEN,
		UCallbackMessage: NotifyIconMessageId,
	}
	nid.CbSize = uint32(unsafe.Sizeof(nid))

	ret, _, _ := Shell_NotifyIcon.Call(NIM_ADD, uintptr(unsafe.Pointer(&nid)))
	if ret == 0 {
		return nil, errors.New("shell notify create failed")
	}

	nid.UVersion = NOTIFYICON_VERSION

	ret, _, _ = Shell_NotifyIcon.Call(NIM_SETVERSION, uintptr(unsafe.Pointer(&nid)))
	if ret == 0 {
		return nil, errors.New("shell notify version failed")
	}

	ni.id = nid.UID
	ni.mhwnd = mhwnd
	ni.hwnd = hwnd

	return ni, nil
}

type _Systray struct {
	iconPath string
	id uint32
	mhwnd uintptr
	hwnd uintptr
	lclick func()
	rclick func()
	dclick func()
}

func NewIconFromFile(filePath string) (uintptr, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return 0, err
	}
	hicon, _, _ := LoadImage.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(absFilePath))),
		IMAGE_ICON,
		0,
		0,
		LR_DEFAULTSIZE|LR_LOADFROMFILE)
	if hicon == 0 {
		return 0, errors.New("load image failed: " + filePath)
	}
	return hicon, nil
}

func RegisterWindow(name string, proc WindowProc) error {
	hinst, _, _ := GetModuleHandle.Call(0)
	if hinst == 0 {
		return errors.New("get module handle failed")
	}
	hicon, _, _ := LoadIcon.Call(0, uintptr(IDI_APPLICATION))
	if hicon == 0 {
		return errors.New("load icon failed")
	}
	hcursor, _, _ := LoadCursor.Call(0, uintptr(IDC_ARROW))
	if hcursor == 0 {
		return errors.New("load cursor failed")
	}

	var wc WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = syscall.NewCallback(proc)
	wc.HInstance = HINSTANCE(hinst)
	wc.HIcon = HICON(hicon)
	wc.HCursor = HCURSOR(hcursor)
	wc.HbrBackground = COLOR_BTNFACE + 1
	wc.LpszClassName = syscall.StringToUTF16Ptr(name)

	atom, _, _ := RegisterClassEx.Call(uintptr(unsafe.Pointer(&wc)))
	if atom == 0 {
		return errors.New("register class failed")
	}
	return nil
}

type WindowProc func(hwnd HWND, msg uint32, wparam, lparam uintptr) uintptr

type NOTIFYICONDATA struct {
	CbSize           uint32
	HWnd             HWND
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            HICON
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GuidItem         GUID
}

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

type WNDCLASSEX struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     HINSTANCE
	HIcon         HICON
	HCursor       HCURSOR
	HbrBackground HBRUSH
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       HICON
}

type MSG struct {
	HWnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type POINT struct {
	X, Y int32
}

type (
	HANDLE uintptr
	HINSTANCE HANDLE
	HCURSOR HANDLE
	HICON HANDLE
	HWND HANDLE
	HGDIOBJ HANDLE
	HBRUSH HGDIOBJ
)

const (
	WM_LBUTTONUP = 0x0202
	WM_LBUTTONDBLCLK = 0x0203
	WM_RBUTTONUP = 0x0205
	WM_USER = 0x0400
	WM_TRAYICON = WM_USER + 69

	WS_EX_APPWINDOW = 0x00040000
	WS_OVERLAPPEDWINDOW = 0X00000000 | 0X00C00000 | 0X00080000 | 0X00040000 | 0X00020000 | 0X00010000
	CW_USEDEFAULT = 0x80000000

	NIM_ADD = 0x00000000
	NIM_MODIFY = 0x00000001
	NIM_DELETE = 0x00000002
	NIM_SETVERSION = 0x00000004

	NIF_MESSAGE = 0x00000001
	NIF_ICON = 0x00000002
	NIF_TIP = 0x00000004
	NIF_STATE = 0x00000008

	NIS_HIDDEN = 0x00000001

	IMAGE_BITMAP = 0
	IMAGE_ICON = 1
	LR_LOADFROMFILE = 0x00000010
	LR_DEFAULTSIZE = 0x00000040

	IDC_ARROW = 32512
	COLOR_WINDOW = 5
	COLOR_BTNFACE = 15

	GWLP_USERDATA = -21
	WS_CLIPSIBLINGS = 0X04000000
	WS_EX_CONTROLPARENT = 0X00010000

	HWND_MESSAGE = ^HWND(2)
	NOTIFYICON_VERSION = 3

	IDI_APPLICATION = 32512
	WM_APP = 32768
	NotifyIconMessageId = WM_APP + iota
)

var (
	kernel32 = syscall.MustLoadDLL("kernel32")
	GetModuleHandle = kernel32.MustFindProc("GetModuleHandleW")
	GetConsoleWindow = kernel32.MustFindProc("GetConsoleWindow")
	GetLastError = kernel32.MustFindProc("GetLastError")

	shell32 = syscall.MustLoadDLL("shell32.dll")
	Shell_NotifyIcon = shell32.MustFindProc("Shell_NotifyIconW")

	user32 = syscall.MustLoadDLL("user32.dll")
	
	GetMessage = user32.MustFindProc("GetMessageW")
	IsDialogMessage = user32.MustFindProc("IsDialogMessageW")
	TranslateMessage = user32.MustFindProc("TranslateMessage")
	DispatchMessage= user32.MustFindProc("DispatchMessageW")
	
	ShowWindow = user32.MustFindProc("ShowWindow")
	UpdateWindow = user32.MustFindProc("UpdateWindow")
	DefWindowProc = user32.MustFindProc("DefWindowProcW")
	RegisterClassEx = user32.MustFindProc("RegisterClassExW")
	GetDesktopWindow = user32.MustFindProc("GetDesktopWindow")
	CreateWindowEx = user32.MustFindProc("CreateWindowExW")

	LoadImage = user32.MustFindProc("LoadImageW")
	LoadIcon = user32.MustFindProc("LoadIconW")
	LoadCursor = user32.MustFindProc("LoadCursorW")
)
