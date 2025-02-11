package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/Cycloctane/dlnapass/internal/upnp"
	"github.com/koron/go-ssdp"
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

	var wg sync.WaitGroup
	var mu sync.Mutex
	i := 0
	for _, d := range devices {
		wg.Add(1)
		go func(d ssdp.Service) {
			defer wg.Done()
			desc, err := upnp.GetDesc(d.Location)
			mu.Lock()
			defer mu.Unlock()

			fmt.Printf(
				"Device #%d: %s\n\t[+] Server: %s\n\t[+] Description URL: %s\n",
				i+1, d.USN, d.Server, d.Location,
			)
			i++
			if err != nil {
				fmt.Print("\t[!] Error: Cannot connect to target device\n\n")
				return
			}
			fmt.Printf(
				"\t[+] DeviceType: %s\n\t[+] ModelName: %s\n\t[+] FriendlyName: %s\n",
				desc.Device.DeviceType, desc.Device.ModelName, desc.Device.FriendlyName,
			)
			for j, s := range desc.Device.ServiceList {
				fmt.Printf("\t[+] Service #%d: %s\n", j+1, s.ServiceType)
			}
			fmt.Print("\n")
		}(d)
	}
	wg.Wait()
}
