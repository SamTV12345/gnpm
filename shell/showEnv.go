package shell

import (
	"fmt"
	"os"
	"strings"

	"github.com/samtv12345/gnpm/filemanagement"
)

func ShowEnv(cwd string) {
	currentParent, err := filemanagement.FindParentLinkModuleDir(cwd)
	if err != nil {
		panic(err)
	}
	currentShell, err := Shell()
	if err != nil {
		panic(err)
	}
	if strings.Contains(currentShell, "cmd.exe") {
		println("set PATH=" + *currentParent + ";%PATH%")
		return
	}

	if strings.Contains(currentShell, "powershell.exe") {
		println("$env:PATH=\"" + *currentParent + ";$env:PATH\"")
		return
	}

	fmt.Printf("export PATH=%s:%s\n", *currentParent, os.Getenv("PATH"))
}
