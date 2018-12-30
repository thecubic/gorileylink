// grl-demo: interacting demo of CC commands (working)

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
	timeout       = flag.Duration("timeout", 10*time.Second, "timeout")
	debug         = flag.Bool("debug", false, "enable debugging messages")
	wg            sync.WaitGroup
	hci           *linux.Device
	ctx           context.Context
	blec          ble.Client
	nameoraddress string
	err           error
	rileylink     *gorileylink.ConnectedRileyLink
	byteme        []byte
	ok            bool
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
		fmt.Println("usage: grl-subs <address-or-name>")
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

	err = rileylink.NotifySubscribe()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("BLE Subscription Failed")
	} else {
		log.Debug("BLE Subscription Successful")
	}

	// BLE methods

	batteryLevel, err = rileylink.GetBatteryLevel()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Battery Level Error")
	} else {
		log.WithFields(log.Fields{
			"batteryLevel": batteryLevel,
		}).Info("Battery Level")
	}

	customName, err = rileylink.GetCustomName()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Custom Name Error")
	} else {
		log.WithFields(log.Fields{
			"customName": customName,
		}).Info("Custom Name")
	}

	bleversion, err = rileylink.GetBLEVersion()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("BLE Version Error")
	} else {
		log.WithFields(log.Fields{
			"bleversion": bleversion,
		}).Info("BLE Version")
	}

	// CC methods

	ok, err = rileylink.GetState()
	if ok {
		log.Info("State: OK")
	} else {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("State: Bad")
	}
	radioversion, err := rileylink.GetRadioVersion()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Radio Version Error")
	} else {
		log.WithFields(log.Fields{
			"radioversion": radioversion,
		}).Info("Radio Version")
	}

	stats, err := rileylink.GetStatistics()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Statistics Error")
	} else {
		log.WithFields(log.Fields{
			"collected":         stats.Collected,
			"uptime":            stats.Uptime,
			"recvoverflows":     stats.RecvOverflows,
			"recvfifooverflows": stats.RecvFifoOverflows,
			"packetsrecv":       stats.PacketsRecv,
			"packetsxmit":       stats.PacketsXmit,
			"crcfails":          stats.CRCFailures,
			"spisyncfails":      stats.SPISyncFailures,
		}).Info("Statistics")
	}

	// This is what procrastinating RF packetry looks like
	log.Info("starting LED dance")
	for _n := 0; _n < 3; _n++ {
		log.WithFields(log.Fields{"green": "on"}).Debug("step")
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDOn)
		log.WithFields(log.Fields{"blue": "on"}).Debug("step + wait")
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDOn)
		time.Sleep(100 * time.Millisecond)

		log.WithFields(log.Fields{"green": "off"}).Debug("step")
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDOff)
		log.WithFields(log.Fields{"blue": "off"}).Debug("step + wait")
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDOff)
		time.Sleep(100 * time.Millisecond)

		log.WithFields(log.Fields{"blue": "on"}).Debug("step")
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDOn)
		log.WithFields(log.Fields{"green": "on"}).Debug("step + wait")
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDOn)
		time.Sleep(100 * time.Millisecond)

		log.WithFields(log.Fields{"green": "off"}).Debug("step")
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDOff)
		log.WithFields(log.Fields{"blue": "off"}).Debug("step + wait")
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDOff)
		time.Sleep(100 * time.Millisecond)

		log.WithFields(log.Fields{"green": "on"}).Debug("step")
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDOn)
		log.WithFields(log.Fields{"blue": "on"}).Debug("step + wait")
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDOn)
		time.Sleep(100 * time.Millisecond)

		log.WithFields(log.Fields{"blue": "off"}).Debug("step")
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDOff)
		log.WithFields(log.Fields{"green": "off"}).Debug("step + wait")
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDOff)
		time.Sleep(100 * time.Millisecond)

		log.WithFields(log.Fields{"blue": "on"}).Debug("step")
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDOn)
		log.WithFields(log.Fields{"green": "on"}).Debug("step + wait")
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDOn)
		time.Sleep(100 * time.Millisecond)

		log.WithFields(log.Fields{"blue": "off"}).Debug("step")
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDOff)
		log.WithFields(log.Fields{"green": "off"}).Debug("step + wait")
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDOff)
		time.Sleep(100 * time.Millisecond)
	}

	// TODO: not working
	// ok, err = rileylink.Reset()
	// if ok {
	// 	fmt.Println("Reset: OK")
	// } else {
	// 	fmt.Printf("Reset: %v\n", err)
	// }

	blec.CancelConnection()
}
