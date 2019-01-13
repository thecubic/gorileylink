// grl-rename: display or change the customizable name of a RileyLink
// e.g. ./grl-rename aa:bb:cc:dd:ee:ff
// e.g. ./grl-rename aa:bb:cc:dd:ee:ff DaveyLink
// e.g. ./grl-rename DaveyLink JimmyLink

package main

import (
	"flag"
	"fmt"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux"
	"github.com/thecubic/gorileylink"
	"golang.org/x/net/context"
)

var (
	timeout          = flag.Duration("timeout", 10*time.Second, "timeout")
	debug            = flag.Bool("debug", false, "enable debugging messages")
	wg               sync.WaitGroup
	hci              *linux.Device
	ctx              context.Context
	blec             ble.Client
	nameoraddress    string
	newfreqi         string
	newfreq          uint64
	err              error
	rileylink        *gorileylink.ConnectedRileyLink
	customNameBefore string
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
		fmt.Println("usage: grl-rename <address-or-name> [new name]")
		return
	}
	newfreqi = flag.Arg(1)
	if newfreqi != "" {
		newfreq, err = strconv.ParseUint(newfreqi, 0, 64)
		if err != nil {
			log.WithFields(log.Fields{
				"input": newfreqi,
				"err":   err,
			}).Fatal("Input Frequency Error")
		}
		if newfreq < 1000000 {
			newfreq = newfreq * 1000000
		}
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

	if newfreqi != "" {
		err = rileylink.SetFrequency(uint32(newfreq))
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"freq": uint32(newfreq),
			}).Fatal("Set Frequency Failed")
		} else {
			log.WithField("frequency", uint32(newfreq)).Debug("Set Frequency")
		}
	}

	frequency, err := rileylink.GetFrequency()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Get Frequency Error")
	} else {
		log.WithFields(log.Fields{
			"frequency": frequency,
		}).Info("Get Frequency")
	}

	// disconnect from rileylink
	blec.CancelConnection()
}
