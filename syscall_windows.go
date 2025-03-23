//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output zsyscall_windows.go syscall_windows.go

package main

//sys AddFontResource(fontPath *uint16) (ret int32, err error) = gdi32.AddFontResourceW
//sys RemoveFontResource(fontPath *uint16) (ret int32, err error) = gdi32.RemoveFontResourceW
//sys SendMessage(hWnd uintptr, msg uint32, wParam uintptr, lParam uintptr) (ret int32, err error) = user32.SendMessageW
//sys GetFontResourceInfo(fontPath *uint16, bufferSize *uint32, buffer uintptr, queryType uint32) (ret int32, err error) = gdi32.GetFontResourceInfoW

const (
	DWINFO_FONT_DESCRIPTION = 1
	DWINFO_FONT_TYPE        = 3
	WM_FONTCHANGE           = 0x001D
	HWND_BROADCAST          = 0xFFFF
)
