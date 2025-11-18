package utils

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func IsMusl() bool {
	if runtime.GOOS != "linux" {
		return false
	}

	// Häufige musl linker-Pfade
	globs := []string{
		"/lib/ld-musl-*.so.1",
		"/lib64/ld-musl-*.so.1",
		"/usr/local/lib/ld-musl-*.so.1",
	}

	for _, g := range globs {
		if matches, _ := filepath.Glob(g); len(matches) > 0 {
			return true
		}
	}

	if out, err := exec.Command("ldd", "--version").CombinedOutput(); err == nil {
		if strings.Contains(string(out), "musl") {
			return true
		}
	}

	return false
}
