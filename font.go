package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func AddFont(fontPath string) error {
	pathPtr, _ := syscall.UTF16PtrFromString(fontPath)
	ret, err := AddFontResource(pathPtr)
	if err != nil {
		if dbg != nil {
			dbg.Error(fmt.Sprintf("AddFont: AddFontResourceW failed: return code=%d, error=%v, winerrno=%d", ret, err, uint32(err.(syscall.Errno))))
		}
		return err
	}
	if dbg != nil {
		dbg.Info(fmt.Sprintf("AddFont: Font loaded successfully: %s, return code=%d", fontPath, ret))
	}
	return nil
}

func RemoveFont(fontPath string) error {
	pathPtr, _ := syscall.UTF16PtrFromString(fontPath)
	ret, err := RemoveFontResource(pathPtr)
	if err != nil {
		if dbg != nil {
			dbg.Error(fmt.Sprintf("RemoveFont: RemoveFontResourceW failed: return code=%d, error=%v", ret, err))
		}
		return err
	}
	if dbg != nil {
		dbg.Info(fmt.Sprintf("RemoveFont: Font unloaded successfully: %s, return code=%d", fontPath, ret))
	}
	return nil
}

func NotifyFontChange() error {
	ret, err := SendMessage(HWND_BROADCAST, WM_FONTCHANGE, 0, 0)
	if err != nil {
		if dbg != nil {
			dbg.Error(fmt.Sprintf("NotifyFontChange: SendMessageW(WM_FONTCHANGE) failed: return code=%d, error=%v", ret, err))
		}
		return err
	}
	if dbg != nil {
		dbg.Info(fmt.Sprintf("NotifyFontChange: WM_FONTCHANGE broadcast sent successfully, return code=%d", ret))
	}
	return nil
}

func GetFontName(fontPath string) (string, error) {
	pathPtr, _ := syscall.UTF16PtrFromString(fontPath)
	var bufferSize uint32

	// First call: Get required buffer size
	ret, err := GetFontResourceInfo(pathPtr, &bufferSize, uintptr(0), DWINFO_FONT_DESCRIPTION)
	if err != nil {
		if dbg != nil {
			dbg.Warn(fmt.Sprintf("GetFontName: GetFontResourceInfoW (first call) failed, font either not loaded or other problem: return code=%d, error=%v, winerrno=%d", ret, err, uint32(err.(syscall.Errno))))
		}
		// we return error here as either the font is not loaded (caller can load it and try again) or other error
		return "", err
	} else {
		if dbg != nil {
			dbg.Info(fmt.Sprintf("GetFontName: GetFontResourceInfoW (first call) successfull: return code=%d, buffersize=%d", ret, bufferSize))
		}
	}

	if bufferSize == 0 {
		if dbg != nil {
			dbg.Error("GetFontResourceInfoW failed: API returned bufferSize = 0, meaning no data is available")
		}
		return "", fmt.Errorf("GetFontName: GetFontResourceInfoW failed: API returned bufferSize = 0, meaning no data is available")
	}

	fontName := make([]uint16, bufferSize/2)

	// Second call: Retrieve actual font name
	ret, err = GetFontResourceInfo(pathPtr, &bufferSize, uintptr(unsafe.Pointer(&fontName[0])), DWINFO_FONT_DESCRIPTION)
	if err != nil {
		if dbg != nil {
			dbg.Error(fmt.Sprintf("GetFontName: GetFontResourceInfoW (second call) failed: return code=%d, error=%v, winerrno=%d", ret, err, uint32(err.(syscall.Errno))))
		}
		return "", err
	}
	fontNameStr := syscall.UTF16ToString(fontName)
	if dbg != nil {
		dbg.Info(fmt.Sprintf("GetFontName: GetFontResourceInfoW (second call) successfull: return code=%d, fontname=%s", ret, fontNameStr))
	}
	return fontNameStr, nil
}

func GetFontNameWithType(fontPath string) (string, error) {
	fontName, err := GetFontName(fontPath)
	if err != nil {
		return "", err
	}
	pathPtr, _ := syscall.UTF16PtrFromString(fontPath)
	var fontType uint32
	bufferSize := uint32(4) // DWORD is always 4 bytes. Right? Right?!
	ret, err := GetFontResourceInfo(pathPtr, &bufferSize, uintptr(unsafe.Pointer(&fontType)), DWINFO_FONT_TYPE)
	if err == nil && fontType >= 1 {
		fontName += " (TrueType)"
	}
	if dbg != nil {
		if err != nil {
			dbg.Error(fmt.Sprintf("GetFontNameWithType: GetFontResourceInfoW: return code=%d, error=%v, winerrno=%d", ret, err, uint32(err.(syscall.Errno))))
		} else {
			dbg.Info(fmt.Sprintf("GetFontNameWithType: GetFontResourceInfoW: return code=%d", ret))
		}
	}
	return fontName, nil
}

func LoadFontFromFile(fontPath string) error {
	if dbg != nil {
		dbg.Info(fmt.Sprintf("Using font file: '%s'", fontPath))
	}
	if info, err := os.Stat(fontPath); err != nil || info.IsDir() {
		if dbg != nil {
			dbg.Error(fmt.Sprintf("PrintFontName: Can't find or open file '%s' , error=%v)", fontPath, err))
		}
		return fmt.Errorf("can't find or open file '%s'", fontPath)
	}

	// Load the font
	if err := AddFont(fontPath); err != nil {
		return err
	}

	// Send WM_FONTCHANGE broadcast
	if err := NotifyFontChange(); err != nil {
		return err
	}

	return nil
}

func UnloadFontFromFile(fontPath string) error {
	if dbg != nil {
		dbg.Info(fmt.Sprintf("Using font file: '%s'", fontPath))
	}
	if info, err := os.Stat(fontPath); err != nil || info.IsDir() {
		if dbg != nil {
			dbg.Error(fmt.Sprintf("PrintFontName: Can't find or open file '%s' , error=%v)", fontPath, err))
		}
		return fmt.Errorf("can't find or open file '%s'", fontPath)
	}

	// Load the font
	if err := RemoveFont(fontPath); err != nil {
		return err
	}

	// Send WM_FONTCHANGE broadcast
	if err := NotifyFontChange(); err != nil {
		return err
	}

	return nil
}

func GetFontNameFromFile(fontPath string) (string, error) {
	if dbg != nil {
		dbg.Info(fmt.Sprintf("Using font file: '%s'", fontPath))
	}
	if info, err := os.Stat(fontPath); err != nil || info.IsDir() {
		if dbg != nil {
			dbg.Error(fmt.Sprintf("PrintFontName: Can't find or open file '%s' , error=%v)", fontPath, err))
		}
		return "", fmt.Errorf("can't find or open file '%s'", fontPath)
	}
	fontIsLoadedBefore := true
	var fontName string
	var err error

	// Retrieve the font name
	fontName, err = GetFontNameWithType(fontPath)
	if err != nil {
		if err == syscall.EINVAL { // font either invalid or not loaded yet
			fontIsLoadedBefore = false

			// Load the font
			if err := AddFont(fontPath); err != nil {
				return "", err
			}

			// Retrieve the font name
			fontName, err = GetFontNameWithType(fontPath)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	if dbg != nil {
		dbg.Info(fmt.Sprintf("font loaded before: %t", fontIsLoadedBefore))
	}
	if !fontIsLoadedBefore {
		// Unload the font, as it hadn't been loaded before we did
		if err := RemoveFont(fontPath); err != nil {
			return "", err
		}
	}
	return fontName, nil

}
