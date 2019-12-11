package sys

import (
	"os"
	"strings"
)

func CnfMustDataDir(dir string) string {
	if dir == "" {
		dir = "~/.hacash_mainnet"
	}
	if strings.HasPrefix(dir, "~/") {
		dir = os.Getenv("HOME") + string([]byte(dir)[1:])
	}
	return dir
}
