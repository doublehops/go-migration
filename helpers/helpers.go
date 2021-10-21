package helpers

import "os"

func PrintMsg(msg string) {
	os.Stderr.WriteString(msg)
}
