package main

import (
	"log"
)

func main() {
	if err := CommandsRoot().Execute(); err != nil {
		log.Fatal(err)
	}
}
