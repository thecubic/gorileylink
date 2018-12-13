package gorileylink

import (
	"fmt"
	"log"

	"github.com/currantlabs/ble"
)

var (
	// from gatt.xml
	rileyLinkSvc  = ble.MustParse("0235733b-99c5-4197-b856-69219c2a3845")
	dataChr       = ble.MustParse("c842e849-5028-42e2-867c-016adada9155")
	respCountChr  = ble.MustParse("6e6c7910-b89e-43a5-a0fe50c5e2b81f4a")
	timerTickChr  = ble.MustParse("6e6c7910-b89e-43a5-78af50c5e2b86f7e")
	customNameChr = ble.MustParse("d93b2af0-1e28-11e4-8c210800200c9a66")
	versionChr    = ble.MustParse("30d99dc9-7c91-4295-a0510a104d238cf2")
	ledModeChr    = ble.MustParse("c6d84241-f1a7-4f9c-a25ffce16732f14e")
)

// ConnectedRileyLink represents a BLE connection to a rileylink
type ConnectedRileyLink struct {
	client       ble.Client
	batterySvc   *ble.Service
	batteryChr   *ble.Characteristic
	rileyLinkSvc *ble.Service
	dataChr      *ble.Characteristic
	// notifier
	respCountChr  *ble.Characteristic
	timerTickChr  *ble.Characteristic
	customNameChr *ble.Characteristic
	versionChr    *ble.Characteristic
	ledModeChr    *ble.Characteristic
}

// on respCountChr notification, dataChr should be read out

// AttachBTLE creates a connection descriptor for a rileylink based on input
// of a legitimate BLE-layer connected device.  It will fail if you give it
// a BT speaker or whatever
func AttachBTLE(blec ble.Client) (*ConnectedRileyLink, error) {
	var (
		err            error
		batterySvcP    *ble.Service
		batteryChrP    *ble.Characteristic
		rileyLinkSvcP  *ble.Service
		dataChrP       *ble.Characteristic
		respCountChrP  *ble.Characteristic
		timerTickChrP  *ble.Characteristic
		customNameChrP *ble.Characteristic
		versionChrP    *ble.Characteristic
		ledModeChrP    *ble.Characteristic
	)
	blep, err := blec.DiscoverProfile(true)
	if err != nil {
		log.Fatalf("couldn't fetch BLE profile")
	}

	for _, s := range blep.Services {
		if s.UUID.Equal(ble.UUID16(0x180F)) {
			batterySvcP = s
			for _, c := range s.Characteristics {
				if c.UUID.Equal(ble.UUID16(0x2a19)) {
					batteryChrP = c
				}
			}
		} else if s.UUID.Equal(rileyLinkSvc) {
			rileyLinkSvcP = s
			for _, c := range s.Characteristics {
				if c.UUID.Equal(dataChr) {
					dataChrP = c
				} else if c.UUID.Equal(respCountChr) {
					respCountChrP = c
				} else if c.UUID.Equal(timerTickChr) {
					timerTickChrP = c
				} else if c.UUID.Equal(customNameChr) {
					customNameChrP = c
				} else if c.UUID.Equal(versionChr) {
					versionChrP = c
				} else if c.UUID.Equal(ledModeChr) {
					ledModeChrP = c
				}
			}
		}
	}

	if batterySvcP == nil {
		return nil, fmt.Errorf("batterySvc missing")
	} else if batteryChrP == nil {
		return nil, fmt.Errorf("batteryChr missing")
	} else if rileyLinkSvcP == nil {
		return nil, fmt.Errorf("rileyLinkSvc missing")
	} else if dataChrP == nil {
		return nil, fmt.Errorf("dataChr missing")
	} else if respCountChrP == nil {
		return nil, fmt.Errorf("respCountChr missing")
	} else if timerTickChrP == nil {
		return nil, fmt.Errorf("timerTickChr missing")
	} else if customNameChrP == nil {
		return nil, fmt.Errorf("customNameChr missing")
	} else if versionChrP == nil {
		return nil, fmt.Errorf("versionChr missing")
	} else if ledModeChrP == nil {
		return nil, fmt.Errorf("ledModeChr missing")
	}

	// yep
	return &ConnectedRileyLink{
		blec,
		batterySvcP,
		batteryChrP,
		rileyLinkSvcP,
		dataChrP,
		respCountChrP,
		timerTickChrP,
		customNameChrP,
		versionChrP,
		ledModeChrP,
	}, nil
}
