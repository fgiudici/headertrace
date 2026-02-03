package main

import (
	"log"

	"github.com/fgiudici/headertrace/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
