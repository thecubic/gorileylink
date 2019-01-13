// grl-rename: display or change the customizable name of a RileyLink
// e.g. ./grl-rename aa:bb:cc:dd:ee:ff
// e.g. ./grl-rename aa:bb:cc:dd:ee:ff DaveyLink
// e.g. ./grl-rename DaveyLink JimmyLink

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
	timeout          = flag.Duration("timeout", 10*time.Second, "timeout")
	debug            = flag.Bool("debug", false, "enable debugging messages")
	wg               sync.WaitGroup
	hci              *linux.Device
	ctx              context.Context
	blec             ble.Client
	nameoraddress    string
	newname          string
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
	newname = flag.Arg(1)

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

	channel := gorileylink.RLPCPump
	rtimeout := 300 * time.Millisecond
	response, err := rileylink.GetPacket(channel, rtimeout)
	if err != nil {
		log.WithFields(log.Fields{
			"channel": channel,
			"timeout": rtimeout,
		}).Error("GetPacket")
	} else {
		log.WithFields(log.Fields{
			"channel": channel,
			"timeout": rtimeout,
			"result":  response.Result,
		}).Info("GetPacket")
	}

	// disconnect from rileylink
	blec.CancelConnection()
}
