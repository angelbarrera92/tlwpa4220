package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/angelbarrera92/tlwpa4220/pkg/tlwpa4220"
)

type cliSubCommand struct {
	flag      *flag.FlagSet
	deviceIP  *string
	username  *string
	password  *string
	command   *string
	tlwpa4220 *tlwpa4220.Client
}

//nolint:revive
func NewCliSubCommand() *cliSubCommand {
	serveflag := flag.NewFlagSet("cli", flag.ExitOnError)

	deviceIP := serveflag.String("device-ip", "", "IP address of the TL-WPA4220 device to monitor")
	username := serveflag.String("username", "admin", "The username required to access the TL-WPA4220 device")
	password := serveflag.String("password", "", "The password required to access the TL-WPA4220 device")

	return &cliSubCommand{
		flag:     serveflag,
		deviceIP: deviceIP,
		username: username,
		password: password,
	}
}

func (f *cliSubCommand) Parse(args []string) error {
	err := f.flag.Parse(args)
	if err != nil {
		return err
	}

	if *f.deviceIP == "" {
		return fmt.Errorf("device-ip is required")
	}

	if *f.password == "" {
		return fmt.Errorf("password is required")
	}

	if f.flag.NArg() < 1 {
		return fmt.Errorf("command is required")
	}

	f.command = &f.flag.Args()[0]

	if *f.command == "" {
		return fmt.Errorf("command is required")
	}

	if !(*f.command == "wireless-statistics" || *f.command == "reboot" || *f.command == "powerline-statistics") {
		return fmt.Errorf("command '%s' is not valid", *f.command)
	}

	f.tlwpa4220 = &tlwpa4220.Client{
		IP:       *f.deviceIP,
		Username: *f.username,
		Password: *f.password,
	}

	return nil
}

func (f *cliSubCommand) Run() error {
	log.Printf("Running command %s", *f.command)

	switch *f.command {
	case "wireless-statistics":
		ws, err := f.tlwpa4220.WirelessStatistics()
		if err != nil {
			log.Printf("Error getting wireless statistics: %s", err)
			return err
		}

		bytes, err := json.Marshal(ws)
		if err != nil {
			log.Printf("Error marshalling wireless statistics: %s", err)
			return err
		}

		log.Printf("Wireless statistics: %v", string(bytes))
	case "powerline-statistics":
		po, err := f.tlwpa4220.PowerLineStatistics()
		if err != nil {
			log.Printf("Error getting powerline statistics: %s", err)
			return err
		}

		bytes, err := json.Marshal(po)
		if err != nil {
			log.Printf("Error marshalling powerline statistics: %s", err)
			return err
		}

		log.Printf("Powerline statistics: %v", string(bytes))
	case "reboot":
		err := f.tlwpa4220.Reboot()
		if err != nil {
			log.Printf("Error rebooting device: %s", err)
			return err
		}
	default:
		return fmt.Errorf("command is not valid")
	}

	return nil
}
