package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

type options struct {
	Networks         []string `short:"n" long:"network" required:"true" value-name:"CIDR" description:"CIDR of the network to be scanned"`
	ExcludedNetworks []string `short:"e" long:"exclude" value-name:"CIDR" description:"CIDR of the excluded networks"`
	ExcludedIPs      []string `long:"exclude-ip" value-name:"IP" description:"Excluded IP from the scan"`
	Execute          string   `short:"x" long:"execute" value-name:"COMMAND" description:"Command to be executed when an host is found" default:"echo {{ .Host }}"`
	Parallel         int      `long:"parallel" description:"Number of goroutine" default:"8"`
	// TODO: Maybe in a future release
	// Verbose          []bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
}

func parseFlags() *options {
	opts := &options{}

	if _, err := flags.Parse(opts); err != nil {
		os.Exit(1)
	}
	return opts
}
