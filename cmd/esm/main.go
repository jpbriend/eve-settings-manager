package main

import (
	"os"

	"github.com/jpbriend/eve-settings-manager/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
