package shell

import (
	"strings"

	"github.com/samtv12345/gnpm/filemanagement"
)

func ShowEnv(cwd string) string {
	currentParent, err := filemanagement.FindParentLinkModuleDir(cwd)
	if err != nil {
		panic(err)
	}
	currentShell, err := Shell()
	if err != nil {
		panic(err)
	}
	if strings.Contains(currentShell, "cmd.exe") {
		return "set PATH=" + *currentParent + ";%PATH%"
	}

	if strings.Contains(currentShell, "powershell.exe") {
		return "$env:PATH=\"" + *currentParent + ";$env:PATH\""
	}

	return "export PATH=" + *currentParent + ":$PATH"
}
