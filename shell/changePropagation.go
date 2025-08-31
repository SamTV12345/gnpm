package shell

import "os"

func PropagateChangesToCurrentShell(path string) error {
	return os.Setenv("PATH", path+string(os.PathListSeparator)+os.Getenv("PATH"))
}
