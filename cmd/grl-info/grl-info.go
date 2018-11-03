package main

import (
	"flag"
	"fmt"
	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux"
	"github.com/thecubic/gorileylink"
	"golang.org/x/net/context"
	"log"
	"sync"
	"time"
)

var (
	timeout   = flag.Duration("timeout", 60*time.Second, "timeout")
	rileylink = flag.String("rileylink", "", "address of the rileylink")
)

func main() {
	flag.Parse()
	if len(*rileylink) == 0 {
		log.Fatalf("must pass rileylink")
	}
	var wg sync.WaitGroup
	d, err := linux.NewDevice()
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)

	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *timeout))

	log.Printf("connecting to %v", *rileylink)
	filter := func(adv ble.Advertisement) bool {
		if len(adv.LocalName()) > 0 {
			if adv.LocalName() == "RileyLink" {
				log.Printf("found a RileyLink: %v", adv.Address().String())
			}
		}
		return adv.Address().String() == *rileylink
	}
	blec, err := ble.Connect(ctx, filter)
	if err != nil {
		log.Fatalf("couldn't connect to %v: %v", rileylink, err)
	} else {
		log.Printf("connected to %v", blec.Address())
	}

	go func() {
		defer wg.Done()
		<-blec.Disconnected()
		log.Printf("disconnected from %v", blec.Address())
	}()
	wg.Add(1)
	defer wg.Wait()

	rl, err := gorileylink.AttachBTLE(blec)
	if err != nil {
		log.Fatalf("couldn't get RileyLink descriptor: %v", err)
	}

	var (
		batteryLevel int
		customName   string
		version      string
		ledMode      gorileylink.LEDMode
	)

	fmt.Printf("Inspecting %v\n", *rileylink)

	customName, err = rl.GetCustomName()
	if err != nil {
		fmt.Printf("couldn't get custom name: %v\n")
	}
	fmt.Printf("  Custom Name: %v\n", customName)

	batteryLevel, err = rl.BatteryLevel()
	if err != nil {
		fmt.Printf("couldn't get battery level: %v\n")
	}
	fmt.Printf("  Battery Level: %v%%\n", batteryLevel)

	version, err = rl.Version()
	if err != nil {
		fmt.Printf("couldn't get version: %v\n")
	}
	fmt.Printf("  Firmware Version: %v\n", version)

	ledMode, err = rl.GetLEDMode()
	if err != nil {
		fmt.Printf("couldn't get LED Mode: %v\n")
	}
	fmt.Printf("  LED Mode: %v\n", ledMode)

	blec.CancelConnection()
}
