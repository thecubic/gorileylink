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
)

func main() {
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

	fmt.Printf("connected: %v\n", nameoraddress)

	err = rileylink.NotifySubscribe()

	fmt.Printf("subscribed -> %v\n", err)
	// err = rileylink.GetVersion()
	// fmt.Printf("GetVersion() -> %v\n", err)
	rileylink.GetStatistics()
	time.Sleep(10 * time.Second)
	blec.CancelConnection()
}
