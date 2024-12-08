package upnp

import (
	"fmt"
	"log"
	"net"
	"runtime"

	"github.com/koron/go-ssdp"
)

const defaultServer = "%s UPnP/%d.%d DLNAPass"

type UpnpAds struct {
	rootAds    [3]*ssdp.Advertiser
	serviceAds []*ssdp.Advertiser
}

func (u *UpnpAds) NotifyDevice() error {
	for _, v := range u.rootAds {
		if err := v.Alive(); err != nil {
			return err
		}
	}
	return nil
}

func (u *UpnpAds) NotifyService() error {
	for _, v := range u.serviceAds {
		if err := v.Alive(); err != nil {
			return err
		}
	}
	return nil
}

func (u *UpnpAds) NotifyAll() error {
	if err := u.NotifyDevice(); err != nil {
		return err
	}
	if err := u.NotifyService(); err != nil {
		return err
	}
	return nil
}

func (u *UpnpAds) CloseAll() {
	for _, v := range u.rootAds {
		v.Bye()
		v.Close()
	}
	for _, v := range u.serviceAds {
		v.Bye()
		v.Close()
	}
}

func SetLogger(logger *log.Logger)    { ssdp.Logger = logger }
func SetInterface(nic *net.Interface) { ssdp.Interfaces = []net.Interface{*nic} }

func advertiserFactory(location, ssdpServer string, defaultMaxAge int) func(st, usn string) (*ssdp.Advertiser, error) {
	return func(st, usn string) (*ssdp.Advertiser, error) {
		return ssdp.Advertise(st, usn, location, ssdpServer, defaultMaxAge)
	}
}

func SetupAdvertise(location string, desc *RootDesc, maxAge int) (*UpnpAds, error) {
	ssdpServer := fmt.Sprintf(defaultServer, runtime.GOOS, desc.SpecVersion.Major, desc.SpecVersion.Minor)
	newAdvertiser := advertiserFactory(location, ssdpServer, maxAge)

	ads := &UpnpAds{}
	ads.rootAds[0], _ = newAdvertiser(
		"upnp:rootdevice",
		fmt.Sprintf("%s::upnp:rootdevice", desc.Device.UDN),
	)
	ads.rootAds[1], _ = newAdvertiser(desc.Device.UDN, desc.Device.UDN)
	ads.rootAds[2], _ = newAdvertiser(
		desc.Device.DeviceType,
		fmt.Sprintf("%s::%s", desc.Device.UDN, desc.Device.DeviceType),
	)
	if err := ads.NotifyDevice(); err != nil {
		return nil, err
	}

	ads.serviceAds = make([]*ssdp.Advertiser, len(desc.Device.ServiceList))
	for k, s := range desc.Device.ServiceList {
		Ad, _ := newAdvertiser(
			s.ServiceType,
			fmt.Sprintf("%s::%s", desc.Device.UDN, s.ServiceType),
		)
		ads.serviceAds[k] = Ad
	}

	if err := ads.NotifyService(); err != nil {
		return nil, err
	}

	return ads, nil
}

func SearchDevice(waitSec int) ([]ssdp.Service, error) {
	service, err := ssdp.Search("upnp:rootdevice", waitSec, "")
	if err != nil {
		return nil, err
	}
	return service, nil
}
