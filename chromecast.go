package chromecast

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/vishen/go-chromecast/application"
)

type Chromecast struct {
	AddrV4 net.IP
	AddrV6 net.IP
	Port   int

	Name string
	Host string

	UUID       string
	Device     string
	Status     string
	DeviceName string
	InfoFields map[string]string
}

// GetUUID returns a unqiue id of a cast entry.
func (c Chromecast) GetUUID() string {
	return c.UUID
}

// GetName returns the identified name of a cast entry.
func (c Chromecast) GetName() string {
	return c.DeviceName
}

// GetAddr returns the IPV4 of a cast entry.
func (c Chromecast) GetAddr() string {
	return fmt.Sprintf("%s", c.AddrV4)
}

// GetPort returns the port of a cast entry.
func (c Chromecast) GetPort() int {
	return c.Port
}

// Discover will return a channel with any cast dns entries found.
func Discover() (<-chan Chromecast, error) {
	var iface *net.Interface

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(3))
	defer cancel()

	var opts = []zeroconf.ClientOption{}
	if iface != nil {
		opts = append(opts, zeroconf.SelectIfaces([]net.Interface{*iface}))
	}

	resolver, err := zeroconf.NewResolver(opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to create new zeroconf resolver: %w", err)
	}

	castDNSEntriesChan := make(chan Chromecast, 5)
	entriesChan := make(chan *zeroconf.ServiceEntry, 5)
	go func() {
		if err := resolver.Browse(ctx, "_googlecast._tcp", "local", entriesChan); err != nil {
			log.Printf("error: unable to browser for mdns entries: %v", err)
			return
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(castDNSEntriesChan)
				return
			case entry := <-entriesChan:
				if entry == nil {
					continue
				}
				chromecast := Chromecast{
					Port: entry.Port,
					Host: entry.HostName,
				}
				if len(entry.AddrIPv4) > 0 {
					chromecast.AddrV4 = entry.AddrIPv4[0]
				}
				if len(entry.AddrIPv6) > 0 {
					chromecast.AddrV6 = entry.AddrIPv6[0]
				}
				infoFields := make(map[string]string, len(entry.Text))
				for _, value := range entry.Text {
					if kv := strings.SplitN(value, "=", 2); len(kv) == 2 {
						key := kv[0]
						val := kv[1]

						infoFields[key] = val

						switch key {
						case "fn":
							chromecast.DeviceName = val
						case "md":
							chromecast.Device = val
						case "id":
							chromecast.UUID = val
						}
					}
				}
				chromecast.InfoFields = infoFields
				castDNSEntriesChan <- chromecast
			}
		}
	}()
	return castDNSEntriesChan, nil
}

func (c Chromecast) Play(file string) error {
	app, err := c.makeApplication()
	if err != nil {
		fmt.Printf("unable to get cast application: %v\n", err)
		return nil
	}

	contentType := ""
	transcode := false
	detach := false

	if err := app.Load(file, contentType, transcode, detach, false); err != nil {
		fmt.Printf("unable to load media: %v\n", err)
		return nil
	}

	return nil
}

func (c Chromecast) makeApplication() (*application.Application, error) {
	applicationOptions := []application.ApplicationOption{
		application.WithDebug(false),
		application.WithCacheDisabled(false),
	}

	app := application.NewApplication(applicationOptions...)
	if err := app.Start(c.GetAddr(), c.GetPort()); err != nil {
		return nil, err
	}
	return app, nil
}
