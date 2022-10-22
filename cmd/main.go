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

type CliCommand interface {
	Parse(args []string) error
	Run() error
}

func printHelp() {
	fmt.Println("Usage: tlwpa4220 [command]")
	fmt.Println("Available commands:")
	fmt.Println("  serve-metrics")
	fmt.Println("  cli")
	// Print also the version
	fmt.Printf("Version: %s, commit: %s, built at: %s\n", version, commit, date)
}

func main() {

	// Print help if no arguments are provided
	if len(os.Args) == 1 {
		printHelp()
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		log.Println("serve-metrics or cli subcommand is required")
		os.Exit(1)
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
