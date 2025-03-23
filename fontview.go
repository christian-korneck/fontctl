package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	fontviewExe = filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "fontview.exe")
)

func PreviewFontWithFontview(fontFile string) error {

	if info, err := os.Stat(fontFile); err != nil || info.IsDir() {
		if dbg != nil {
			dbg.Error(fmt.Sprintf("PrintFontName: Can't find or open file '%s' , error=%v)", fontFile, err))
		}
		return fmt.Errorf("can't find or open file '%s'", fontFile)
	}

	cmd := exec.Command(fontviewExe, fontFile)
	err := cmd.Start()
	if err != nil {
		return (err)
	}
	return nil
}
