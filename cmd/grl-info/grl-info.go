// grl-info: display name, BLE FW version, battery level of RileyLink
// e.g. ./grl-info aa:bb:cc:dd:ee:ff
// e.g. ./grl-info DaveyLink

package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux"
	"github.com/thecubic/gorileylink"
	"golang.org/x/net/context"
)

var (
	timeout       = flag.Duration("timeout", 10*time.Second, "timeout")
	wg            sync.WaitGroup
	hci           *linux.Device
	ctx           context.Context
	blec          ble.Client
	nameoraddress string
	err           error
	rileylink     *gorileylink.ConnectedRileyLink
	batteryLevel  int
	rssi          int
	customName    string
	version       string
)

func main() {
	flag.Parse()
	nameoraddress = flag.Arg(0)
	if nameoraddress == "" {
		fmt.Println("usage: grl-info <address-or-name>")
		return
	}

	// boilerplate connect to rileylink
	hci, ctx = gorileylink.OpenBLE(*timeout)
	blec, err = gorileylink.ConnectNameOrAddress(ctx, nameoraddress)
	if err != nil {
		log.Fatalf("couldn't connect to %v: %v", nameoraddress, err)
	}

	rileylink, err = gorileylink.AttachBTLE(blec)
	if err != nil {
		log.Fatalf("couldn't bind %v as RileyLink: %v", nameoraddress, err)
	}

	// launch a goroutine to wrap BLE disconnection for a clean exit
	go func() {
		defer wg.Done()
		<-blec.Disconnected()
	}()
	wg.Add(1)
	defer wg.Wait()
	// end boilerplate connect to rileylink

	batteryLevel, err = rileylink.BatteryLevel()
	if err != nil {
		fmt.Printf("couldn't get battery level: %v\n", err)
	}

	customName, err = rileylink.GetCustomName()
	if err != nil {
		fmt.Printf("couldn't get custom name: %v\n", err)
	}

	version, err = rileylink.Version()
	if err != nil {
		fmt.Printf("couldn't get version: %v\n", err)
	}

	rssi = rileylink.ReadRSSI()

	fmt.Printf("%v @ %v (%v dBm): %v %v %v%%\n", nameoraddress, blec.Address().String(), rssi, customName, version, batteryLevel)

	// disconnect from rileylink
	blec.CancelConnection()
}
