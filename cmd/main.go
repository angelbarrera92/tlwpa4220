package main

import (
	"flag"
	"log"
	"os"

	"github.com/angelbarrera92/tlwpa4220/cmd/cli"
	"github.com/angelbarrera92/tlwpa4220/cmd/serve"
)

type CliCommand interface {
	Parse(args []string) error
	Run() error
}

func main() {

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
