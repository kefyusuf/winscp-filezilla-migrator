package main

import (
	"fmt"
	"os"

	"github.com/muety/winscp2filezilla/app"
)

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}