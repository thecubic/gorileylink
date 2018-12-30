// grl-info: display name, BLE FW version, battery level of RileyLink
// e.g. ./grl-info aa:bb:cc:dd:ee:ff
// e.g. ./grl-info DaveyLink

package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux"
	"github.com/thecubic/gorileylink"
	"golang.org/x/net/context"
)

var (
	timeout       = flag.Duration("timeout", 10*time.Second, "connection timeout")
	debug         = flag.Bool("debug", false, "enable debugging messages")
	wg            sync.WaitGroup
	hci           *linux.Device
	ctx           context.Context
	blec          ble.Client
	nameoraddress string
	err           error
	rileylink     *gorileylink.ConnectedRileyLink
	batteryLevel  int
	customName    string
	bleversion    string
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
		fmt.Println("usage: grl-info <address-or-name>")
		os.Exit(1)
	}

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

	// BLE methods

	batteryLevel, err = rileylink.GetBatteryLevel()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Battery Level Error")
	} else {
		log.WithFields(log.Fields{
			"batteryLevel": batteryLevel,
		}).Debug("Battery Level")
	}

	customName, err = rileylink.GetCustomName()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Custom Name Error")
	} else {
		log.WithFields(log.Fields{
			"customName": customName,
		}).Debug("Custom Name")
	}

	bleversion, err = rileylink.GetBLEVersion()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("LE Version Error")
	} else {
		log.WithFields(log.Fields{
			"bleversion": bleversion,
		}).Debug("BLE Version")
	}

	fmt.Printf("%v @ %v: %v %v %v%%\n", nameoraddress, blec.Address().String(), customName, bleversion, batteryLevel)

	// disconnect from rileylink
	blec.CancelConnection()
}
