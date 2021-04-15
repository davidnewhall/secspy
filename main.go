package main

import (
	"log"
	"os"
	"time"

	"github.com/davidnewhall/secspy/cli"
)

// Keep it simple.
func main() {
	// Set time zone based on TZ env variable.
	setTimeZone(os.Getenv("TZ"))

	if err := cli.Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}

func setTimeZone(tz string) {
	if tz == "" {
		return
	}

	var err error

	if time.Local, err = time.LoadLocation(tz); err != nil {
		log.Printf("[ERROR] Loading TZ Location '%s': %v", tz, err)
	}
}
