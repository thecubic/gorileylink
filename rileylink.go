package gorileylink

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

func (crl *ConnectedRileyLink) SetCustomName(newname string) error {
	var (
		data []byte
		err  error
	)
	data = []byte(newname)
	err = crl.client.WriteCharacteristic(crl.customNameChr, data, false)
	return err
}

func (crl *ConnectedRileyLink) SetLEDMode(mode LEDMode) error {
	var (
		err error
	)
	err = crl.client.WriteCharacteristic(crl.ledModeChr, []byte{byte(mode)}, false)
	return err
}

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
