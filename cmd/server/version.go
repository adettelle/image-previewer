package main

import (
	"fmt"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func GetVersion() string {
	return fmt.Sprintf("release: %s - buildDate: %s - gitHash: %s\n", release, buildDate, gitHash)
}
