// grl-leds: display or change diagnostic LEDs on RileyLink
// e.g. ./grl-leds aa:bb:cc:dd:ee:ff on
// e.g. ./grl-leds DaveyLink on
// e.g. ./grl-leds DaveyLink

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
	timeout         = flag.Duration("timeout", 10*time.Second, "timeout")
	wg              sync.WaitGroup
	hci             *linux.Device
	ctx             context.Context
	blec            ble.Client
	nameoraddress   string
	err             error
	rileylink       *gorileylink.ConnectedRileyLink
	desiredledstate string
)

func main() {
	flag.Parse()
	nameoraddress = flag.Arg(0)
	if nameoraddress == "" {
		fmt.Println("usage: grl-leds <address-or-name> [off/on/auto]")
		return
	}
	desiredledstate = flag.Arg(1)

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

	if desiredledstate == "on" {
		err = rileylink.SetLEDMode(gorileylink.LEDModeOn)
	} else if desiredledstate == "off" {
		err = rileylink.SetLEDMode(gorileylink.LEDModeOff)
	} else if desiredledstate == "auto" {
		err = rileylink.SetLEDMode(gorileylink.LEDModeAuto)
	}
	if err != nil {
		fmt.Printf("error in setting LED Mode: %v\n", err)
	}

	ledmode, err := rileylink.GetLEDMode()
	if err != nil {
		fmt.Printf("couldn't get LED Mode: %v\n", err)
	}
	if ledmode == gorileylink.LEDModeOn {
		fmt.Printf("%v: LED Mode: on\n", nameoraddress)
	} else if ledmode == gorileylink.LEDModeOff {
		fmt.Printf("%v: LED Mode: off\n", nameoraddress)
	} else if ledmode == gorileylink.LEDModeAuto {
		fmt.Printf("%v: LED Mode: auto\n", nameoraddress)
	}

	// disconnect from rileylink
	blec.CancelConnection()
}
