package main

import (
	"log"

	"github.com/davidnewhall/secspy/cli"
)

func main() {
	err := cli.Start()
	if err != nil {
		log.Fatal(err)
	}
}
