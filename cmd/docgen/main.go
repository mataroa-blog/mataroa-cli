package main

import (
	"log"

	"git.sr.ht/~glorifiedgluer/mata/commands"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmd := commands.SetupCommands()

	header := &doc.GenManHeader{
		Title:   "mata",
		Section: "1",
		Manual:  "General Commands Manual",
		// TODO: provide a proper date here
		// Date:    &dt,
		Source: "https://git.sr.ht/~glorifiedgluer/mata",
	}

	err := doc.GenManTree(cmd, header, "doc/result")
	if err != nil {
		log.Fatal(err)
	}
}
