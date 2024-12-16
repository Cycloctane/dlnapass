package main

import (
	"flag"
	"fmt"

	"github.com/Cycloctane/dlnapass/internal/upnp"
)

const defaultSearchSec = 1

func main() {
	searchSec := flag.Uint("t", defaultSearchSec, "Search duration (seconds)")
	flag.Parse()

	fmt.Println("Searching UPnP root devices...")
	devices, err := upnp.SearchDevice(int(*searchSec))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d devices:\n\n", len(devices))
	for i, d := range devices {
		fmt.Printf("Device #%d: %s\n", i+1, d.USN)
		fmt.Printf("\t[+] Server: %s\n", d.Server)
		fmt.Printf("\t[+] Description URL: %s\n", d.Location)
		desc, err := upnp.GetDesc(d.Location)
		if err != nil {
			fmt.Print("\t[!] Error: Cannot connect to target device\n\n")
			continue
		}
		fmt.Printf("\t[+] DeviceType: %s\n", desc.Device.DeviceType)
		fmt.Printf("\t[+] ModelName: %s\n", desc.Device.ModelName)
		fmt.Printf("\t[+] FriendlyName: %s\n", desc.Device.FriendlyName)
		for j, s := range desc.Device.ServiceList {
			fmt.Printf("\t[+] Service #%d: %s\n", j+1, s.ServiceType)
		}
		fmt.Print("\n")
	}
}
