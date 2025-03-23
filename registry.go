package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func CreateWindowsFontRegistryKey(fontName, fontFile string, user bool) error {
	if !user { // For HLKM use only the filename
		fontFile = filepath.Base(fontFile)
	}

	var baseKey registry.Key
	if user {
		baseKey = registry.CURRENT_USER
	} else {
		baseKey = registry.LOCAL_MACHINE
	}

	fontsKeyPath := `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`

	k, err := registry.OpenKey(baseKey, fontsKeyPath, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %v", err)
	}
	defer k.Close()

	names, err := k.ReadValueNames(0)
	if err != nil {
		return fmt.Errorf("failed to read registry value names: %v", err)
	}

	for _, name := range names {
		val, _, err := k.GetStringValue(name)
		if err == nil && strings.EqualFold(val, fontFile) {
			if dbg != nil {
				dbg.Warn(fmt.Sprintf("CreateWindowsFontRegistryKey: Font file '%s' is already registered under key '%s'.\n", fontFile, name))
			}
			return nil
		}
	}

	newFontName := fontName
	for _, name := range names {
		if strings.EqualFold(name, fontName) {
			existingFile, _, err := k.GetStringValue(name)
			if err == nil {
				// If the font name exists but points to the same file, no action is needed.
				if strings.EqualFold(existingFile, fontFile) {
					return nil
				}
				// Otherwise, choose a new key name with an incrementing suffix.
				index := 1
				for {
					candidate := fmt.Sprintf("%s (%d)", fontName, index)
					exists := false
					for _, n := range names {
						if strings.EqualFold(n, candidate) {
							exists = true
							break
						}
					}
					if !exists {
						newFontName = candidate
						if dbg != nil {
							dbg.Warn(fmt.Sprintf("CreateWindowsFontRegistryKey: Font name '%s' already exists with a different file. Using new key name '%s'.\n", fontName, newFontName))
						}
						break
					}
					index++
				}
				break
			}
		}
	}

	if err := k.SetStringValue(newFontName, fontFile); err != nil {
		return fmt.Errorf("failed to set registry value: %v", err)
	}

	return nil
}

func RemoveWindowsFontRegistryKeys(fontFile string, user bool) error {
	if !user { // for HKLM use only the filename
		fontFile = filepath.Base(fontFile)
	}

	var baseKey registry.Key
	if user {
		baseKey = registry.CURRENT_USER
	} else {
		baseKey = registry.LOCAL_MACHINE
	}

	fontsKeyPath := `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`

	k, err := registry.OpenKey(baseKey, fontsKeyPath, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %v", err)
	}
	defer k.Close()

	names, err := k.ReadValueNames(0)
	if err != nil {
		return fmt.Errorf("failed to read registry value names: %v", err)
	}

	found := false
	for _, name := range names {
		val, _, err := k.GetStringValue(name)
		if err == nil && strings.EqualFold(val, fontFile) {
			if err := k.DeleteValue(name); err != nil {
				return fmt.Errorf("failed to delete registry value '%s': %v", name, err)
			}
			if dbg != nil {
				dbg.Warn(fmt.Sprintf("Deleted registry key '%s' that pointed to '%s'\n", name, fontFile))
			}
			found = true
		}
	}

	if !found {
		if dbg != nil {
			dbg.Warn(fmt.Sprintf("RemoveWindowsFontRegistryKeys: No registry keys found for font file '%s'\n", fontFile))
		}

	}

	return nil
}
