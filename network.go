package main

import (
	"fmt"
	"net"
	"os"

	funk "github.com/thoas/go-funk"
)

func nextIP(ip net.IP, bit int) {
	ip[bit]++
	if ip[bit] == 0 {
		nextIP(ip, bit-1)
	}
}

func calcIPs(network string) []string {
	var ips []string

	ip, ipnet, err := net.ParseCIDR(network)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Network '%s' skipped: %v\n", network, err)
		return ips
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); nextIP(ip, len(ip)-1) {
		ips = append(ips, ip.String())
	}

	// Remove Network + Broadcast
	return ips[1 : len(ips)-1]
}

func getIPs(opts *options) []string {
	var ips, allIPs, excludedIPs []string

	// Generate all IPs from networks
	for _, net := range opts.Networks {
		allIPs = append(allIPs, calcIPs(net)...)
	}

	// Generate all excluded IPs from networks and IPs
	for _, net := range opts.ExcludedNetworks {
		excludedIPs = append(excludedIPs, calcIPs(net)...)
	}
	for _, ip := range opts.ExcludedIPs {
		if net.ParseIP(ip) == nil {
			fmt.Fprintf(os.Stderr, "Excluded IP '%s' skipped: unvalid IP\n", ip)
		} else {
			excludedIPs = append(excludedIPs, ip)
		}
	}

	// Exclude IPs from all IPs
	ips = funk.Chain(allIPs).Uniq().Filter(func(ip string) bool {
		return !funk.Contains(excludedIPs, allIPs)
	}).Value().([]string)

	return ips
}
