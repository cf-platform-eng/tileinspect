package version

import "fmt"

type VersionOpt struct {
}

var Version = "dev"

func (_ *VersionOpt) Execute(args []string) error {
	fmt.Printf("tileinspect version: %s\n", Version)
	return nil
}