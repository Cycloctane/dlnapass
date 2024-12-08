# DLNAPass

DLNA Pass: Pass DLNA/UPnP discovery messages through subnet/vpn.

This program acts like a [SSDP Server Proxy](https://datatracker.ietf.org/doc/html/draft-cai-ssdp-v1-01#section-7.2) for UPnP devices. It retrieves description from UPnP device and sends SSDP messages in local network to make the target UPnP device visible in local subnet.

Location header in advertisement remains the same as UPnP device's orginal address. UPnP clients (control points) should be able to connect to the target UPnP device directly.

It can also be used for accessing remote DLNA/UPnP service through vpns that do not route multicast traffic (like ipsec and openvpn).

## Usage

### DLNAPass

```bash
./dlnapass -i $interface -u $description_url -t 1800
```

- `-i`: Network interface used for sending SSDP multicast messages. Note that linux loopback interfaces may need extra configuration to support multicast.
- `-u`: URL of target UPnP device's root device description. (`http://host:8200/rootDesc.xml` for minidlna)
- `-t`: Advertisement duration (max-age).
- `--verbose`: Enable go-ssdp logger to see more verbose logging.

Use control-c to send `ssdp:byebye` and exit the program gracefully.

For example, to make a minidlna server at 192.168.1.2 in remote network visible after connecting to the remote network with ipsec:

```bash
./dlnapass -u http://192.168.1.2:8200/rootDesc.xml
```

DLNAPass will send `ssdp:alive` and respond to `ssdp:discover`. DLNA clients on localhost can be notified with server's original address (192.168.1.2) and connect to server through vpn tunnel.

### DLNAFind

DLNAFind is a helper program for searching available UPnP root devices and DLNA (UPnP-AV) MediaServer service in local network. It is useful for getting root device description URL that can be used in DLNAPass `-u`.

```bash
./dlnafind
```
