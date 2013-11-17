package systray

import (
	"errors"
	"path/filepath"
	"unsafe"
	"syscall"
)

func (p *_Systray) Run() error {
	proc := func(hwnd HWND, msg uint32, wparam, lparam uintptr) uintptr {
		if lparam == WM_LBUTTONDBLCLK {
			p.dclick()
		} else if (lparam == WM_LBUTTONUP) {
			p.lclick()
		} else if (lparam ==  WM_RBUTTONUP) {
			p.rclick()
		}
		result, _, _ := DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wparam, lparam)
		return result
	}

	name := syscall.StringToUTF16Ptr("Systray")

	hinst, _, _ := GetModuleHandle.Call(0)
	if (hinst == 0) {
		return errors.New("can't get module handle")
	}
	cursor, _, _ := LoadCursor.Call(0, IDC_ARROW)
	if (cursor == 0) {
		return errors.New("can't get cursor")
	}
	icon, _, _ := LoadIcon.Call(0, uintptr(IDI_APPLICATION))
	if icon == 0 {
		return errors.New("can't get app icon")
	}

	var wclass WNDCLASSEX
	wclass.CbSize = uint32(unsafe.Sizeof(wclass))
	wclass.LpfnWndProc = syscall.NewCallback(proc)
	wclass.HInstance = HINSTANCE(hinst)
	wclass.LpszClassName = name
	wclass.HbrBackground = COLOR_WINDOW + 1
	wclass.HCursor = HCURSOR(cursor)
	wclass.HIcon = HICON(icon)

	result, _, _ := RegisterClassEx.Call(uintptr(unsafe.Pointer(&wclass)))
	if (result == 0) {
		return errors.New("reg win class failed")
	}
	desktop, _, _ := GetDesktopWindow.Call()
	if (desktop == 0) {
		return errors.New("get desktop failed")
	}

	p.hwin, _, _ = CreateWindowEx.Call(
		WS_EX_APPWINDOW,
		uintptr(unsafe.Pointer(name)),
		uintptr(unsafe.Pointer(name)),
		WS_OVERLAPPEDWINDOW ,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		desktop,
		0,
		hinst,
		0)

	if (p.hwin == 0) {
		return errors.New("create win failed")
	}

	msg := uintptr(unsafe.Pointer(&MSG{}))
	for {
		result, _, _ := GetMessage.Call(msg, 0, 0, 0)
		if (result == 0) {
			break
		}
		TranslateMessage.Call(msg);
		DispatchMessage.Call(msg);
	}
	return nil
}

func (p *_Systray) Stop() error {
	if (p.nid == nil) {
		return nil
	}
	result, _, _ := Shell_NotifyIcon.Call(NIM_DELETE, uintptr(unsafe.Pointer(p.nid)))
	if (result == 0) {
		return errors.New("hide icon failed")
	}
	return nil
}

func (p *_Systray) Show(file string, hint string) error {
	path := filepath.Join(p.iconPath, file)
	icon, _, _ := LoadImage.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		IMAGE_ICON,
		0,
		0,
		LR_LOADFROMFILE | LR_DEFAULTSIZE)

	if (icon == 0) {
		icon, _, _ = LoadIcon.Call(0, IDI_APPLICATION)
	}
	if (icon == 0) {
		return errors.New("can't load icon")
	}

	flag := uint32(NIF_ICON | NIF_MESSAGE)
	if len(hint) != 0 {
		flag = flag | NIF_TIP
	}

	var nid *NOTIFYICONDATA
	if (p.nid != nil) {
		nid = p.nid
	} else {
		nid = &NOTIFYICONDATA{}
		nid.CbSize = uint32(unsafe.Sizeof(nid))
		nid.HWnd = HWND(p.hwin)
		nid.UID = 1
		nid.UCallbackMessage = WM_TRAYICON
	}

	nid.UFlags = flag
	nid.HIcon = HICON(icon)
	copy(nid.SzTip[:], syscall.StringToUTF16(hint))

	var result uintptr
	if (p.nid != nil) {
		result, _, _ = Shell_NotifyIcon.Call(NIM_MODIFY, uintptr(unsafe.Pointer(nid)))
	} else {
		result, _, _ = Shell_NotifyIcon.Call(NIM_ADD, uintptr(unsafe.Pointer(nid)))
		p.nid = nid
	}
	if (result == 0) {
		return errors.New("shell show call failed")
	}
	return nil
}

func (p *_Systray) OnClick(fun func()) {
	p.lclick = fun
	p.rclick = fun
	p.dclick = fun
}

func _NewSystray(iconPath string, clientPath string, port int) _Systray {
	return _Systray{iconPath, 0, nil, func(){}, func(){}, func(){}}
}

type _Systray struct {
	iconPath string
	hwin uintptr
	nid *NOTIFYICONDATA
	lclick func()
	rclick func()
	dclick func()
}

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

	NIF_MESSAGE = 0x00000001
	NIF_ICON = 0x00000002
	NIF_TIP = 0x00000004

	IMAGE_BITMAP = 0
	IMAGE_ICON = 1
	LR_LOADFROMFILE = 0x00000010
	LR_DEFAULTSIZE = 0x00000040

	IDC_ARROW = 32512
	COLOR_WINDOW = 5

	IDI_APPLICATION = 32512
)

var (
	kernel32 = syscall.MustLoadDLL("kernel32")
	GetModuleHandle = kernel32.MustFindProc("GetModuleHandleW")
	GetConsoleWindow = kernel32.MustFindProc("GetConsoleWindow")

	shell32 = syscall.MustLoadDLL("shell32.dll")
	Shell_NotifyIcon = shell32.MustFindProc("Shell_NotifyIconW")

	user32 = syscall.MustLoadDLL("user32.dll")
	GetMessage = user32.MustFindProc("GetMessageW")
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
