package main

import (
	"fmt"
	"os"

	"github.com/kefyusuf/winscp-filezilla-migrator/app"
)

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
