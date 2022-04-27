package main

import (
	"log"
	"time"

	"git.sr.ht/~glorifiedgluer/mata/commands"
	"git.sr.ht/~glorifiedgluer/mata/mataroa"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmd := commands.SetupCommands()

	// This date represents the day I discovered how to make this
	// program truly reproducible
	dt, _ := time.Parse(mataroa.ISO8601Layout, "2022-04-27")

	header := &doc.GenManHeader{
		Title:   "mata",
		Section: "1",
		Manual:  "General Commands Manual",
		Date:    &dt,
		Source:  "https://git.sr.ht/~glorifiedgluer/mata",
	}

	err := doc.GenManTree(cmd, header, "doc/result")
	if err != nil {
		log.Fatal(err)
	}
}
