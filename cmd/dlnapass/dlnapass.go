package main

import (
	"flag"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/Cycloctane/dlnapass/internal/upnp"
)

const defaultMaxAgeSec = 1800

func main() {
	locationStr := flag.String("u", "", "URL of upnp device's root description xml")
	nicStr := flag.String("i", "", "Network interface for multicast")
	maxAge := flag.Int("t", defaultMaxAgeSec, "Max age of upnp notify in seconds")
	verbose := flag.Bool("verbose", false, "Verbose logging (enable go-ssdp logger)")
	flag.Parse()
	logger := log.New(os.Stdout, "", 0)
	if *maxAge < defaultMaxAgeSec {
		logger.Fatalln("Error: Max-age should be greater than 1800 seconds.")
	}
	location, err := url.Parse(*locationStr)
	if err != nil || !location.IsAbs() {
		logger.Fatalln("Error: Invalid root description URL.")
	}

	if *nicStr != "" {
		nic, err := net.InterfaceByName(*nicStr)
		if err != nil {
			logger.Panicln(err)
		}
		upnp.SetInterface(nic)
	}

	desc, err := upnp.GetDesc(location.String())
	if err != nil {
		logger.Fatalf("Cannot get root description with provided URL: %s\n", err)
	}

	if *verbose {
		upnp.SetLogger(logger)
	}
	logger.Printf(
		"DeviceType: %s\nModelName: %s\nFriendlyName: %s\n\n",
		desc.Device.DeviceType, desc.Device.ModelName, desc.Device.FriendlyName,
	)
	logger.SetFlags(log.Ltime)
	logger.SetPrefix("[+] ")
	logger.Println("Starting SSDP advertisement...")

	ads, err := upnp.SetupAdvertise(location.String(), desc, defaultMaxAgeSec)
	if err != nil {
		logger.Panicln(err)
	}

	repeat := time.Tick(time.Duration(defaultMaxAgeSec) * time.Second)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

loop:
	for {
		select {
		case <-quit:
			break loop
		case <-repeat:
			if !upnp.IsAlive(location.String(), desc.Device.UDN) {
				logger.Printf("Cannot reach original device at %s. Closing...\n", location)
				break loop
			} else {
				if err := ads.NotifyAll(); err != nil {
					break loop
				}
			}
		}
	}

	ads.CloseAll()
}
