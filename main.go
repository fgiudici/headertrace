package main

import (
	"github.com/fgiudici/headertrace/cmd"
	"github.com/fgiudici/headertrace/pkg/logging"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logging.Fatalf("%v", err)
	}
}
