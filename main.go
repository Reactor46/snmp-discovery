package main

import (
	"bytes"
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
	wg := sync.WaitGroup{}

	for _, ips := range chunks {
		wg.Add(1)
		go func(ips []string, wg *sync.WaitGroup) {
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
							fmt.Fprintf(os.Stderr, "Failed to create a pinger: %v\n", err)
							continue
						}
					}
				} else {
					if pinger, err = ping.New("0.0.0.0", ""); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to create a pinger: %v\n", err)
						continue
					}
				}

				_, err = pinger.Ping(addr, 2500*time.Millisecond)
				pinger.Close()
				if err != nil {
					continue
				}

				host := ip
				names, err := net.LookupAddr(ip)
				if err == nil && len(names) > 0 {
					host = names[0]
				}

				var buffer bytes.Buffer
				command.Execute(&buffer, struct{ Host string }{host})
				exec.Command("/bin/sh", "-c", buffer.String()).Run()
				fmt.Printf("[%s] Execute '%s'\n", ip, buffer.String())
			}
		}(ips, &wg)
	}

	wg.Wait()
}
