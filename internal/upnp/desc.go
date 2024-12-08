package upnp

import (
	"encoding/xml"
	"net/http"
)

type RootDesc struct {
	XMLName     xml.Name `xml:"root"`
	SpecVersion struct {
		Major int `xml:"major"`
		Minor int `xml:"minor"`
	} `xml:"specVersion"`
	Device *struct {
		UDN          string
		DeviceType   string `xml:"deviceType"`
		FriendlyName string `xml:"friendlyName"`
		ModelName    string `xml:"modelName"`
		ServiceList  []*struct {
			ServiceType string `xml:"serviceType"`
		} `xml:"serviceList>service"`
	} `xml:"device"`
}

func GetDesc(location string) (*RootDesc, error) {
	resp, err := http.Get(location)
	if err != nil {
		return nil, err
	}
	desc := &RootDesc{}
	decoder := xml.NewDecoder(resp.Body)
	if err := decoder.Decode(desc); err != nil {
		return nil, err
	}
	resp.Body.Close()
	return desc, nil
}

func IsAlive(location, lastUDN string) bool {
	desc, err := GetDesc(location)
	if err != nil || desc.Device.UDN != lastUDN {
		return false
	}
	return true
}
