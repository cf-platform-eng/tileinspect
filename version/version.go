package version

import "fmt"

type VersionOpt struct {
}

var Version = "dev"

func (_ *VersionOpt) Execute(args []string) error {
	fmt.Printf("marman version: %s\n", Version)
	return nil
}