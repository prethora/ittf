package commands

import (
	"fmt"

	"github.com/prethora/ittf/internal/version"
)

func execVersion() {
	fmt.Println(version.Version)
}
