// Package main is the entrypoint for the tlwpa4220 exporter.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/angelbarrera92/tlwpa4220/cmd/cli"
	"github.com/angelbarrera92/tlwpa4220/cmd/serve"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// CliCommand is the interface for subcommands.
type CliCommand interface {
	Parse(args []string) error
	Run() error
}

func printHelp() {
	fmt.Printf("Usage: tlwpa4220 [command]\nAvailable commands:\n  serve-metrics\n  cli\nVersion: %s, commit: %s, built at: %s\n", version, commit, date)
}

func main() {
	if len(os.Args) == 1 {
		printHelp()
		os.Exit(0)
	}

	var cmd CliCommand

	switch os.Args[1] {
	case "serve-metrics":
		cmd = serve.NewServeMetricsSubCommand()
	case "cli":
		cmd = cli.NewCliSubCommand()
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := cmd.Parse(os.Args[2:])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = cmd.Run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
