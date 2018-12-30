// grl-leds: display or change diagnostic LEDs on RileyLink
// e.g. ./grl-leds aa:bb:cc:dd:ee:ff on
// e.g. ./grl-leds DaveyLink on
// e.g. ./grl-leds DaveyLink

package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux"
	"github.com/thecubic/gorileylink"
	"golang.org/x/net/context"
)

var (
	timeout         = flag.Duration("timeout", 10*time.Second, "timeout")
	debug           = flag.Bool("debug", false, "enable debugging messages")
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

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

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
		log.WithFields(log.Fields{
			"rileylink": nameoraddress,
			"err":       err,
		}).Fatal("connection failed")
	} else {
		log.WithFields(log.Fields{
			"rileylink": nameoraddress,
		}).Debug("connection succeeded")
	}

	rileylink, err = gorileylink.AttachBTLE(blec)
	if err != nil {
		log.WithFields(log.Fields{
			"rileylink": nameoraddress,
			"err":       err,
		}).Fatal("couldn't bind connected device as RileyLink")
	} else {
		log.WithFields(log.Fields{
			"rileylink": nameoraddress,
		}).Debug("bind as RileyLink succeeded")
	}

	// launch a goroutine to wrap BLE disconnection for a clean exit
	go func() {
		defer wg.Done()
		<-blec.Disconnected()
	}()
	wg.Add(1)
	// this will delay program exit until cleanly disconnected.
	// since this is probably Bluetooth-API-over-IPC, not doing
	// this may persist undesired HCI state
	defer wg.Wait()
	// end boilerplate connect to rileylink

	if desiredledstate == "on" {
		err = rileylink.SetLEDMode(gorileylink.LEDOn)
	} else if desiredledstate == "off" {
		err = rileylink.SetLEDMode(gorileylink.LEDOff)
	} else if desiredledstate == "auto" {
		err = rileylink.SetLEDMode(gorileylink.LEDAuto)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"rileylink": nameoraddress,
			"err":       err,
			"desired":   desiredledstate,
		}).Fatal("Set LED Mode Error")
	}

	ledmode, err := rileylink.GetLEDMode()
	if err != nil {
		log.WithFields(log.Fields{
			"rileylink": nameoraddress,
			"err":       err,
		}).Fatal("Get LED Mode Error")
	}

	if ledmode == gorileylink.LEDOn {
		log.WithFields(log.Fields{
			"rileylink": nameoraddress,
			"leds":      "on",
		}).Info("LED Mode")
	} else if ledmode == gorileylink.LEDOff {
		log.WithFields(log.Fields{
			"rileylink": nameoraddress,
			"leds":      "off",
		}).Info("LED Mode")
	} else if ledmode == gorileylink.LEDAuto {
		log.WithFields(log.Fields{
			"rileylink": nameoraddress,
			"leds":      "auto",
		}).Info("LED Mode")
	}

	// disconnect from rileylink
	blec.CancelConnection()
}
