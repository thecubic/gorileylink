// rileylink.go contains the application-layer specifics of the device

package gorileylink

import (
	"encoding/binary"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/currantlabs/ble"
)

// ConnectedRileyLink represents a BLE connection to a rileylink
type ConnectedRileyLink struct {
	client              ble.Client
	batterySvc          *ble.Service
	batteryChr          *ble.Characteristic
	rileyLinkSvc        *ble.Service
	dataChr             *ble.Characteristic
	respCountChr        *ble.Characteristic
	respCountClientDesc *ble.Descriptor
	timerTickChr        *ble.Characteristic
	customNameChr       *ble.Characteristic
	versionChr          *ble.Characteristic
	ledModeChr          *ble.Characteristic
	rawResponse         chan []byte
	response            chan RLCCResponse
	notifier            chan int
}

// AttachBTLE creates a connection descriptor for a rileylink based on input
// of a legitimate BLE-layer connected device.  It will fail if you give it
// a BT speaker or whatever
// Effectively the constructor
func AttachBTLE(blec ble.Client) (*ConnectedRileyLink, error) {
	var (
		err                  error
		batterySvcP          *ble.Service
		batteryChrP          *ble.Characteristic
		rileyLinkSvcP        *ble.Service
		dataChrP             *ble.Characteristic
		respCountChrP        *ble.Characteristic
		respCountClientDescP *ble.Descriptor
		timerTickChrP        *ble.Characteristic
		customNameChrP       *ble.Characteristic
		versionChrP          *ble.Characteristic
		ledModeChrP          *ble.Characteristic
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
					for _, d := range c.Descriptors {
						if d.UUID.Equal(ble.UUID16(0x2902)) {
							respCountClientDescP = d
						}
					}
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
		respCountClientDescP,
		timerTickChrP,
		customNameChrP,
		versionChrP,
		ledModeChrP,
		make(chan []byte),
		make(chan RLCCResponse),
		make(chan int),
	}, nil
}

// func (crl *ConnectedRileyLink)

// there are two RPC layers to a RileyLink; the BLE113 and the CC1110
// connected to it via SPI.  These run different firmwares; the BLE
// firmware is "ble_rfspy" and is effectively the supervising chip
// and is interacted with directly via GATT characteristics
// The CC1110 firmware is "subg_rfspy" and is interacted with over a
// GATT call-and-response scheme

// ReadRSSI [local] just exposes the underlying call
func (crl *ConnectedRileyLink) ReadRSSI() int {
	return crl.client.ReadRSSI()
}

// NotifySubscribe [local] wires a function as a callback to the data notifier
func (crl *ConnectedRileyLink) NotifySubscribe() error {
	var (
		err error
	)
	// prepare ourselves for notification
	err = crl.client.Subscribe(crl.respCountChr, false, crl.notifyRespCallback)
	if err != nil {
		log.Fatalf("local subscribe failed: %s", err)
	}
	// tell the device to notify us, m'kay
	err = crl.client.WriteDescriptor(crl.respCountClientDesc, enableNotificationValue)
	if err != nil {
		log.Fatalf("remote notify failed: %s", err)
	}
	return err
}

// notifyRespCallback is a simple callback to convert BLE notification
// events into application channel notification events
// such that a reciever knows it now should read data
func (crl *ConnectedRileyLink) notifyRespCallback(dumpval []byte) {
	// NOTE: reading the data characteristic here does not work
	log.WithField("sequence", int(dumpval[0])).Debug("RespCount notified")
	crl.notifier <- int(dumpval[0])
}

// BatteryLevel [BLE] retrieves an approximated battery percentage from the device
func (crl *ConnectedRileyLink) GetBatteryLevel() (int, error) {
	var (
		data []byte
		err  error
	)
	data, err = crl.client.ReadCharacteristic(crl.batteryChr)
	if err != nil {
		return -1, err
	}
	return int(data[0]), nil
}

// GetCustomName [BLE] returns the device's name set by the user
func (crl *ConnectedRileyLink) GetCustomName() (string, error) {
	var (
		data []byte
		err  error
	)
	data, err = crl.client.ReadCharacteristic(crl.customNameChr)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetCustomName [BLE] pushes a new name to the device
func (crl *ConnectedRileyLink) SetCustomName(newname string) error {
	var (
		data []byte
		err  error
	)
	data = []byte(newname)
	err = crl.client.WriteCharacteristic(crl.customNameChr, data, false)
	return err
}

// GetLEDMode [BLE] retrieves the mode of the diagnostic LEDs (blue)
func (crl *ConnectedRileyLink) GetLEDMode() (LEDMode, error) {
	var (
		err  error
		mode LEDMode
		data []byte
	)
	data, err = crl.client.ReadCharacteristic(crl.ledModeChr)
	if err == nil {
		mode = LEDMode(data[0])
	}
	return mode, err
}

// SetLEDMode [BLE] switches the mode of the diagnostic LEDs (blue)
func (crl *ConnectedRileyLink) SetLEDMode(mode LEDMode) error {
	var (
		err error
	)
	err = crl.client.WriteCharacteristic(crl.ledModeChr, []byte{byte(mode)}, false)
	return err
}

// Version [BLE] returns the BLE firmware revision on the device
func (crl *ConnectedRileyLink) GetBLEVersion() (string, error) {
	var (
		data []byte
		err  error
	)
	data, err = crl.client.ReadCharacteristic(crl.versionChr)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// RLCCResponse represents a return from the CC chip
type RLCCResponse struct {
	Result  RileyLinkCCResponseType
	Payload []byte
	RSSI    int
}

// writeCCPacket [CC] pushes an application packet to the CC chip
func (crl *ConnectedRileyLink) writeCCPacket(packet []byte) error {
	lenpluspacket := make([]byte, len(packet)+1)
	lenpluspacket[0] = byte(len(packet))
	copy(lenpluspacket[1:], packet)
	log.WithField("packet", lenpluspacket).Debug("writeCCPacket")
	return crl.client.WriteCharacteristic(crl.dataChr, lenpluspacket, false)
}

// resetCC is just a conveience function for the oneway
func (crl *ConnectedRileyLink) resetCC() error {
	return crl.writeCCPacket([]byte{byte(RLCReset)})
}

// commandCC is just a convenience function for wrapping CC commands
func (crl *ConnectedRileyLink) commandCC(cmd RileyLinkCommand) (*RLCCResponse, error) {
	err := crl.writeCCPacket([]byte{byte(cmd)})
	if err != nil {
		return nil, err
	}
	select {
	case <-crl.notifier:
		log.Debug("notifier fired")
		return crl.readResponse()
	case <-time.After(1 * time.Second):
		log.Debug("notifier did not fire")
		return crl.readResponse()
	}
}

// payloadCommandCC sends a command with an extensible payload
func (crl *ConnectedRileyLink) payloadCommandCC(cmd RileyLinkCommand, payload []byte) (*RLCCResponse, error) {
	fullpacket := make([]byte, len(payload)+1)
	fullpacket[0] = byte(cmd)
	copy(fullpacket[1:], payload)
	err := crl.writeCCPacket(fullpacket)
	if err != nil {
		log.WithField("err", err).Error("writeCCPacket Error")
		return nil, err
	}
	log.Debug("waiting on notifier")

	select {
	case <-crl.notifier:
		log.Debug("notifier fired")
		return crl.readResponse()
	case <-time.After(1 * time.Second):
		log.Debug("notifier did not fire")
		return crl.readResponse()
	}
}

// instantPayloadCommandCC is for when a command with params
// won't trigger a characteristic notification
// so don't wait for what won't happen
func (crl *ConnectedRileyLink) instantPayloadCommandCC(cmd RileyLinkCommand, payload []byte) (*RLCCResponse, error) {
	fullpacket := make([]byte, len(payload)+1)
	fullpacket[0] = byte(cmd)
	copy(fullpacket[1:], payload)
	err := crl.writeCCPacket(fullpacket)
	if err != nil {
		log.WithField("err", err).Error("writeCCPacket Error")
		return nil, err
	}
	return crl.readResponse()
}

func (crl *ConnectedRileyLink) readResponse() (*RLCCResponse, error) {
	var (
		respPayload []byte
		err         error
	)
	responded := false

	log.Debug("readResponse")
	for !responded {
		respPayload, err := crl.client.ReadCharacteristic(crl.dataChr)
		if err != nil {
			log.WithField("err", err).Error("ReadCharacteristic Error")
			return nil, err
		} else {
			log.WithField("payload", respPayload).Debug("ReadCharacteristic")
		}
		if len(respPayload) > 0 {
			responded = true
			log.WithField("lenPayload", len(respPayload)).Debug("captured response")
		} else {
			log.Debug("no response yet, retrying")
		}
	}
	log.WithFields(log.Fields{
		"len": len(respPayload),
	}).Debug("pre-create")
	response := &RLCCResponse{
		RileyLinkCCResponseType(respPayload[0]),
		make([]byte, len(respPayload)-1),
		crl.client.ReadRSSI()}
	log.Debug("pre-copy")
	copy(response.Payload, respPayload[1:])
	log.Debug("returning response")
	return response, err
}

// see https://github.com/ps2/subg_rfspy/blob/master/protocol.md

// Interrupt [CC] is, like, stop what you're doing
func (crl *ConnectedRileyLink) Interrupt() error {
	var err error
	_, err = crl.commandCC(RLCInterrupt)
	return err
}

// GetState [CC] is an internal diagnostic call
func (crl *ConnectedRileyLink) GetState() (bool, error) {
	response, err := crl.commandCC(RLCGetState)
	if err != nil {
		return false, err
	} else if response.Result != RLRSuccess {
		return false, fmt.Errorf("Bad result: %v", response.Result)
	} else if string(response.Payload) != "OK" {
		return false, fmt.Errorf("Not OK: %v", string(response.Payload))
	}
	return true, err
}

// GetRadioVersion [CC] returns the version of the CC firmware
func (crl *ConnectedRileyLink) GetRadioVersion() (string, error) {
	response, err := crl.commandCC(RLCGetVersion)
	if response.Result != RLRSuccess {
		return "", fmt.Errorf("Bad result: %v", response.Result)
	}
	return string(response.Payload), err
}

// GetPacket [CC] does a thing that will be documented at some point
func (crl *ConnectedRileyLink) GetPacket(rlpc RileyLinkPacketChannel, timeout time.Duration) (*RLCCResponse, error) {
	payload := make([]byte, 5)
	payload[0] = byte(rlpc)
	binary.BigEndian.PutUint32(payload[1:], uint32(timeout/time.Millisecond))
	response, err := crl.payloadCommandCC(RLCGetPacket, payload)
	if err != nil {
		log.WithFields(log.Fields{
			"timeout": timeout,
			"channel": rlpc,
			"payload": payload,
			"err":     err,
		}).Error("GetPacket")
		return nil, err
	} else {
		log.WithFields(log.Fields{
			"timeout": timeout,
			"channel": rlpc,
			"payload": payload,
		}).Debug("GetPacket")
	}
	return response, err
}

// SendPacket [CC] does a thing that will be documented at some point
func (crl *ConnectedRileyLink) SendPacket() error {
	var err error
	_, err = crl.commandCC(RLCSendPacket)
	return err
}

// SendAndListen [CC] does a thing that will be documented at some point
func (crl *ConnectedRileyLink) SendAndListen() error {
	var err error
	_, err = crl.commandCC(RLCSendAndListen)
	return err
}

// UpdateRegister [CC] does a thing that will be documented at some point
func (crl *ConnectedRileyLink) UpdateRegister() error {
	var err error
	_, err = crl.commandCC(RLCUpdateRegister)
	return err
}

// RawReset [CC] resets the CC chip and that's that
func (crl *ConnectedRileyLink) RawReset() error {
	return crl.resetCC()
}

// Reset [CC] resets the CC chip, and returns a state call after 100ms
func (crl *ConnectedRileyLink) Reset() (bool, error) {
	var err error
	err = crl.RawReset()
	if err != nil {
		return false, err
	}
	// wait for 100ms for the CC to reset
	time.Sleep(100 * time.Millisecond)
	return crl.GetState()
}

// LED [CC] does a thing that will be documented at some point
func (crl *ConnectedRileyLink) LED(ledc LEDColor, ledm LEDMode) error {
	response, err := crl.payloadCommandCC(RLCLED, []byte{byte(ledc), byte(ledm)})
	if response.Result != RLRSuccess {
		return fmt.Errorf("Bad result: %v", response.Result)
	}
	return err
}

// ReadRegister [CC] reads a cute 'lil 8-bit register
func (crl *ConnectedRileyLink) ReadRegister(reg CxRegister) (byte, error) {
	var (
		err      error
		response *RLCCResponse
	)
	log.Debug("reading register")
	if true {
		// TODO: firmware revision gate
		// subg_rfspy versions < 2.3 need to be told twice
		response, err = crl.instantPayloadCommandCC(RLCReadRegister, []byte{byte(reg), byte(reg)})
	} else {
		response, err = crl.instantPayloadCommandCC(RLCReadRegister, []byte{byte(reg)})
	}
	log.Debug("read register")
	if response.Result != RLRSuccess {
		return 0, fmt.Errorf("Bad result: %v", response.Result)
	}
	return response.Payload[0], err
}

// WriteRegister [CC] writes a cute 'lil 8-bit register
func (crl *ConnectedRileyLink) WriteRegister(reg CxRegister, value byte) error {
	log.WithFields(log.Fields{
		"register": reg,
		"value":    value,
	}).Debug("WriteRegister")
	response, err := crl.instantPayloadCommandCC(RLCUpdateRegister, []byte{byte(reg), value})
	if response.Result != RLRSuccess {
		return fmt.Errorf("Bad result: %v", response.Result)
	}
	return err
}

// SetModeRegisters [CC] does a thing that will be documented at some point
func (crl *ConnectedRileyLink) SetModeRegisters() error {
	var err error
	_, err = crl.commandCC(RLCSetModeRegisters)
	return err
}

// SetSWEncoding [CC] does a thing that will be documented at some point
func (crl *ConnectedRileyLink) SetSWEncoding() error {
	var err error
	_, err = crl.commandCC(RLCSetSWEncoding)
	return err
}

// SetPreamble [CC] does a thing that will be documented at some point
func (crl *ConnectedRileyLink) SetPreamble() error {
	var err error
	_, err = crl.commandCC(RLCSetPreamble)
	return err
}

// ResetRadioConfig [CC] does a thing that will be documented at some point
func (crl *ConnectedRileyLink) ResetRadioConfig() error {
	var err error
	_, err = crl.commandCC(RLCResetRadioConfig)

	return err
}

// GetStatistics [CC] does a thing that will be documented at some point
func (crl *ConnectedRileyLink) GetStatistics() (*RileyLinkStatistics, error) {
	response, err := crl.commandCC(RLCGetStatistics)
	if err != nil {
		return nil, err
	} else if response.Result != RLRSuccess {
		return nil, fmt.Errorf("Bad result: %v", response.Result)
	}
	return &RileyLinkStatistics{
		time.Now(),
		time.Duration(binary.BigEndian.Uint32(response.Payload[0:4])) * time.Millisecond,
		binary.BigEndian.Uint16(response.Payload[4:6]),
		binary.BigEndian.Uint16(response.Payload[6:8]),
		binary.BigEndian.Uint16(response.Payload[8:10]),
		binary.BigEndian.Uint16(response.Payload[10:12]),
		binary.BigEndian.Uint16(response.Payload[12:14]),
		binary.BigEndian.Uint16(response.Payload[14:16]),
	}, err
}
