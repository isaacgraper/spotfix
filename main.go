package main

import (
	"log"

	cmd "github.com/isaacgraper/spotfix.git/internal/cmd/cli"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Printf("Error while trying to run bot: %v", err)
	}
}
