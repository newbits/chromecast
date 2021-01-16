package cmd

import (
	"context"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	castdns "github.com/vishen/go-chromecast/dns"
)

// Run the LS command
func Run() {
	ifaceName := ""
	dnsTimeoutSeconds := 3

	var iface *net.Interface
	var err error

	if ifaceName != "" {
		if iface, err = net.InterfaceByName(ifaceName); err != nil {
			log.Fatalf("unable to find interface %q: %v", ifaceName, err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(dnsTimeoutSeconds))

	defer cancel()

	castEntryChan, err := castdns.DiscoverCastDNSEntries(ctx, iface)
	i := 1

	for d := range castEntryChan {
		fmt.Printf("%d) device=%q device_name=%q address=\"%s:%d\" uuid=%q\n", i, d.Device, d.DeviceName, d.AddrV4, d.Port, d.UUID)
		i++
	}

	if i == 1 {
		fmt.Printf("no cast devices found on network\n")
	}
}
