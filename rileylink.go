package gorileylink

import (
	"fmt"
	"log"
)

// BatteryLevel retrieves an approximated battery percentage from the device
func (crl *ConnectedRileyLink) BatteryLevel() (int, error) {
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

// GetCustomName returns the device's name set by the user
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

// SetCustomName pushes a new name to the device
func (crl *ConnectedRileyLink) SetCustomName(newname string) error {
	var (
		data []byte
		err  error
	)
	data = []byte(newname)
	err = crl.client.WriteCharacteristic(crl.customNameChr, data, false)
	return err
}

// GetLEDMode retrieves the mode of the diagnostic LEDs (blue)
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

// SetLEDMode switches the mode of the diagnostic LEDs (blue)
func (crl *ConnectedRileyLink) SetLEDMode(mode LEDMode) error {
	var (
		err error
	)
	err = crl.client.WriteCharacteristic(crl.ledModeChr, []byte{byte(mode)}, false)
	return err
}

// Version returns the BLE firmware revision on the device
func (crl *ConnectedRileyLink) Version() (string, error) {
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

func (crl *ConnectedRileyLink) GetState() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCGetState)}, true)
	return err
}

func (crl *ConnectedRileyLink) GetVersion() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCGetVersion)}, true)
	return err
}

func (crl *ConnectedRileyLink) GetPacket() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCGetPacket)}, true)
	return err
}

func (crl *ConnectedRileyLink) SendPacket() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCSendPacket)}, true)
	return err
}

func (crl *ConnectedRileyLink) SendAndListen() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCSendAndListen)}, true)
	return err
}

func (crl *ConnectedRileyLink) UpdateRegister() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCUpdateRegister)}, true)
	return err
}

func (crl *ConnectedRileyLink) Reset() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCReset)}, true)
	return err
}

func (crl *ConnectedRileyLink) LED() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCLED)}, true)
	return err
}

func (crl *ConnectedRileyLink) ReadRegister() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCReadRegister)}, true)
	return err
}

func (crl *ConnectedRileyLink) SetModeRegisters() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCSetModeRegisters)}, true)
	return err
}

func (crl *ConnectedRileyLink) SetSWEncoding() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCSetSWEncoding)}, true)
	return err
}

func (crl *ConnectedRileyLink) SetPreamble() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCSetPreamble)}, true)
	return err
}

func (crl *ConnectedRileyLink) ResetRadioConfig() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCResetRadioConfig)}, true)
	return err
}

func (crl *ConnectedRileyLink) GetStatistics() error {
	var err error
	err = crl.client.WriteCharacteristic(crl.dataChr, []byte{byte(RLCGetStatistics)}, true)
	return err
}

// NotifySubscribe wires a function as a callback to the data notifier
func (crl *ConnectedRileyLink) NotifySubscribe() error {
	var err error
	if err = crl.client.Subscribe(crl.respCountChr, true, crl.notifyRespCallback); err != nil {
		log.Fatalf("subscribe failed: %s", err)
	}
	//if err = crl.client.Subscribe(crl.dataChr, false, crl.notifyDataCallback); err != nil {
	//	log.Fatalf("subscribe failed: %s", err)
	//}
	return err
}

func (crl *ConnectedRileyLink) notifyRespCallback(req []byte) {
	fmt.Printf("notifiedResp: %v\n", req)
}

func (crl *ConnectedRileyLink) notifyDataCallback(req []byte) {
	fmt.Printf("notifiedData: %v\n", req)
}
