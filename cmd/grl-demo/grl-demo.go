// grl-subs: WIP: buildup of subscription-based GATT RPC channel
// nope, it does not work yet

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
	byteme        []byte
)

func main() {
	var ok bool
	flag.Parse()
	nameoraddress = flag.Arg(0)
	if nameoraddress == "" {
		fmt.Println("usage: grl-subs <address-or-name>")
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

	err = rileylink.NotifySubscribe()

	ok, err = rileylink.GetState()
	if ok {
		fmt.Println("State: OK")
	} else {
		fmt.Printf("State Bad: %v\n", err)
	}
	radioversion, err := rileylink.GetRadioVersion()
	if err != nil {
		fmt.Printf("Radio Version Err: %v\n", err)
	} else {
		fmt.Printf("Radio Version: %v\n", radioversion)
	}

	stats, err := rileylink.GetStatistics()
	if err != nil {
		fmt.Printf("GetStatistics Err: %v\n", err)
	} else {
		fmt.Printf("Statistics:\n"+
			"  Uptime: %v\n"+
			"  RecvOverflows: %v\n"+
			"  RecvFifoOverlows: %v\n"+
			"  PacketsRecv: %v\n"+
			"  PacketsXmit: %v\n"+
			"  CRCFailures: %v\n"+
			"  SPISyncFailures: %v\n",
			stats.Uptime,
			stats.RecvOverflows,
			stats.RecvFifoOverflows,
			stats.PacketsRecv,
			stats.PacketsXmit,
			stats.CRCFailures,
			stats.SPISyncFailures)
	}

	for _n := 0; _n < 3; _n++ {
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDModeOn)
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDModeOn)
		time.Sleep(100 * time.Millisecond)
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDModeOff)
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDModeOff)
		time.Sleep(100 * time.Millisecond)

		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDModeOn)
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDModeOn)
		time.Sleep(100 * time.Millisecond)
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDModeOff)
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDModeOff)
		time.Sleep(100 * time.Millisecond)

		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDModeOn)
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDModeOn)
		time.Sleep(100 * time.Millisecond)
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDModeOff)
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDModeOff)
		time.Sleep(100 * time.Millisecond)

		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDModeOn)
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDModeOn)
		time.Sleep(100 * time.Millisecond)
		rileylink.LED(gorileylink.LEDBlue, gorileylink.LEDModeOff)
		rileylink.LED(gorileylink.LEDGreen, gorileylink.LEDModeOff)
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
