// grl-rename: display or change the customizable name of a RileyLink
// e.g. ./grl-rename aa:bb:cc:dd:ee:ff
// e.g. ./grl-rename aa:bb:cc:dd:ee:ff DaveyLink
// e.g. ./grl-rename DaveyLink JimmyLink

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
	timeout          = flag.Duration("timeout", 10*time.Second, "timeout")
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

	customNameBefore, err = rileylink.GetCustomName()
	if err != nil {
		fmt.Printf("couldn't get custom name: %v\n", err)
	}

	if len(newname) == 0 {
		fmt.Printf("%v: named %v\n", nameoraddress, customNameBefore)
	} else {
		err = rileylink.SetCustomName(newname)
		if err != nil {
			fmt.Printf("error in renaming: %v\n", err)
		}
		fmt.Printf("%v: renamed %v to %v\n", nameoraddress, customNameBefore, newname)
	}

	// disconnect from rileylink
	blec.CancelConnection()
}
