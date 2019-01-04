package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
	"text/template"
	"time"

	ping "github.com/digineo/go-ping"
	funk "github.com/thoas/go-funk"
)

func main() {
	opts := parseFlags()

	ips := getIPs(opts)
	command := template.Must(template.New("command").Parse(opts.Execute))

	chunks := funk.Chunk(ips, len(ips)/opts.Parallel+1).([][]string)
	wg := &sync.WaitGroup{}

	for _, ips := range chunks {
		go func(ips []string) {
			wg.Add(1)
			defer wg.Done()

			for _, ip := range ips {
				var addr *net.IPAddr
				var err error
				var pinger *ping.Pinger

				if addr, err = net.ResolveIPAddr("ip4", ip); err != nil {
					if addr, err = net.ResolveIPAddr("ip6", ip); err != nil {
						continue
					} else {
						if pinger, err = ping.New("", "::"); err != nil {
							fmt.Fprintf(os.Stderr, "Failed to create a pinger: %v", err)
							continue
						}
					}
				} else {
					if pinger, err = ping.New("0.0.0.0", ""); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to create a pinger: %v", err)
						continue
					}
				}

				_, err = pinger.Ping(addr, 2500*time.Millisecond)
				pinger.Close()
				if err != nil {
					continue
				}

				var buffer bytes.Buffer
				command.Execute(&buffer, struct{ Host string }{ip})
				exec.CommandContext(context.Background(), "/bin/sh", "-c", buffer.String()).Run()
				fmt.Printf("[%s] Execute '%s'\n", ip, buffer.String())
			}
		}(ips)
	}
	wg.Wait()
}
