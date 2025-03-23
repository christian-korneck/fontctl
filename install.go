package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func InstallFontFromFile(fontPath string, installSystemWide bool) error {

	if dbg != nil {
		dbg.Info(fmt.Sprintf("Using font file: '%s'", fontPath))
	}

	if info, err := os.Stat(fontPath); err != nil || info.IsDir() {
		if dbg != nil {
			dbg.Error(fmt.Sprintf("PrintFontName: Can't find or open file '%s' , error=%v)", fontPath, err))
		}
		return fmt.Errorf("can't find or open file '%s'", fontPath)
	}

	var destPath string
	switch installSystemWide {
	case true:
		winDir, err := windows.GetSystemWindowsDirectory()
		if err != nil {
			return fmt.Errorf("can't find Windows system dir")
		}
		destPath = filepath.Join(winDir, "Fonts")
	case false:
		localAppData, err := windows.KnownFolderPath(windows.FOLDERID_LocalAppData, 0)
		if err != nil {
			return fmt.Errorf("can't find User localappdata dir")
		}
		destPath = filepath.Join(localAppData, "Microsoft", "Windows", "Fonts")
	}

	if dbg != nil {
		dbg.Info(fmt.Sprintf("Using destination font dir '%s'", destPath))
	}

	if fi, err := os.Stat(destPath); err != nil {
		if os.IsNotExist(err) {
			if !installSystemWide {
				// if user font dir does not exist, make it
				err = os.MkdirAll(destPath, 0755)
				if err != nil {
					return fmt.Errorf("user Font dir '%s' does not exist and trying to create it failed (%s)", destPath, err)
				}
			} else {
				// abort if system font dir does not exist
				return fmt.Errorf("system font dir '%s' does not exist", destPath)
			}
		}
	} else if !fi.IsDir() {
		return fmt.Errorf("windows font dir path '%s' exists but is not a directory", destPath)
	}

	if dbg != nil {
		dbg.Info(fmt.Sprintf("destination font dir '%s' exists and can be used", destPath))
	}

	err := CopyFile(fontPath, destPath, false)
	if err != nil {
		return err
	}

	fontDestPath := filepath.Join(destPath, filepath.Base(fontPath))

	if fi, err := os.Stat(fontDestPath); err != nil || fi.IsDir() {
		if dbg != nil {
			dbg.Error(fmt.Sprintf("InstallFontFromFile: Can't find or open destination font file after copy '%s' , error=%v)", fontPath, err))
		}
		return fmt.Errorf("can't find or open file '%s' at destination path after copy", fontDestPath)
	}

	// Load the font
	if err := AddFont(fontDestPath); err != nil {
		return err
	}

	// Retrieve the font name
	fontName, err := GetFontNameWithType(fontDestPath)
	if err != nil {
		return err

	}

	err = CreateWindowsFontRegistryKey(fontName, fontDestPath, !installSystemWide)
	if err != nil {
		return err
	}

	// Send WM_FONTCHANGE broadcast
	if err := NotifyFontChange(); err != nil {
		return err
	}

	return nil

}

func UninstallFontFromFile(fontPath string, removeSystemWide bool) error {

	if dbg != nil {
		dbg.Info(fmt.Sprintf("Using font file: '%s'", fontPath))
	}

	if info, err := os.Stat(fontPath); err != nil || info.IsDir() {
		if dbg != nil {
			dbg.Error(fmt.Sprintf("PrintFontName: Can't find or open file '%s' , error=%v)", fontPath, err))
		}
		return fmt.Errorf("can't find or open file '%s'", fontPath)
	}

	var destPath string
	switch removeSystemWide {
	case true:
		winDir, err := windows.GetSystemWindowsDirectory()
		if err != nil {
			return fmt.Errorf("can't find Windows system dir")
		}
		destPath = filepath.Join(winDir, "Fonts")
	case false:
		localAppData, err := windows.KnownFolderPath(windows.FOLDERID_LocalAppData, 0)
		if err != nil {
			return fmt.Errorf("can't find User localappdata dir")
		}
		destPath = filepath.Join(localAppData, "Microsoft", "Windows", "Fonts")
	}

	if dbg != nil {
		dbg.Info(fmt.Sprintf("Using destination font dir '%s'", destPath))
	}

	if fi, err := os.Stat(destPath); err != nil {
		if os.IsNotExist(err) {
			if !removeSystemWide {
				return fmt.Errorf("user Font dir '%s' does not exist and trying to create it failed (%s)", destPath, err)
			}
		} else {
			// abort if system font dir does not exist
			return fmt.Errorf("system font dir '%s' does not exist", destPath)
		}
	} else if !fi.IsDir() {
		return fmt.Errorf("windows font dir path '%s' exists but is not a directory", destPath)
	}

	if dbg != nil {
		dbg.Info(fmt.Sprintf("destination font dir '%s' exists and can be used", destPath))
	}

	fontDestPath := filepath.Join(destPath, filepath.Base(fontPath))

	/* srcHash, err := hashFile(fontPath)
	if err != nil {
		return fmt.Errorf("could not hash source file '%s' (%s), uninstall aborted.", fontPath, err)
	}
	dstHash, err := hashFile(fontDestPath)
	if err != nil {
		return fmt.Errorf("could not hash dest file '%s' (%s), uninstall aborted.", fontDestPath, err)
	}

	if string(srcHash) != string(dstHash) {
		return fmt.Errorf("file '%s' is not idential to file '%s' (hash), uninstall aborted.", fontPath, fontDestPath)
	} */

	err := UnloadFontFromFile(fontDestPath)
	if err != nil {
		if dbg != nil {
			dbg.Warn(fmt.Sprintf("Failed to unload font from file during uninstall. This can be okay. fontfile=%s, error=%s", fontDestPath, err))
		}
	}

	err = RemoveWindowsFontRegistryKeys(fontDestPath, !removeSystemWide)
	if err != nil {
		if dbg != nil {
			dbg.Warn(fmt.Sprintf("Failed finding and removing font registry key during uninstall. This can be okay. fontfile=%s, error=%s", fontDestPath, err))
		}
	}

	err = os.Remove(fontDestPath)
	if err != nil {
		return fmt.Errorf("failed to remove file '%s' (%s), uninstall incomplete", fontPath, err)
	}

	// Step 2: Send WM_FONTCHANGE broadcast
	if err := NotifyFontChange(); err != nil {
		if dbg != nil {
			dbg.Warn(fmt.Sprintf("failed to send WM_FONTCHANGE broadcast during uninstall. This can be okay. fontfile=%s, error=%s", fontDestPath, err))
		}
	}

	return nil

}
