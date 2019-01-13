// rfspy.go contains the CC-layer code for making RF okay again

package gorileylink

import (
	log "github.com/sirupsen/logrus"
)

type RxFilter byte

type SwEncoding byte

type CxRegister byte

const (
	RxFilterWide       RxFilter   = 0x50 // 300KHz
	RxFilterNarrow     RxFilter   = 0x90 // 150KHz
	EncodingNone       SwEncoding = 0x00
	EncodingManchester SwEncoding = 0x01
	Encoding4b6b       SwEncoding = 0x02
	RegisterSync1      CxRegister = 0x00
	RegisterSync0      CxRegister = 0x01
	RegisterPktlen     CxRegister = 0x02
	RegisterPktctrl1   CxRegister = 0x03
	RegisterPktctrl0   CxRegister = 0x04
	RegisterFsctrl1    CxRegister = 0x07
	RegisterFreq2      CxRegister = 0x09
	RegisterFreq1      CxRegister = 0x0a
	RegisterFreq0      CxRegister = 0x0b
	RegisterMdmcfg4    CxRegister = 0x0c
	RegisterMdmcfg3    CxRegister = 0x0d
	RegisterMdmcfg2    CxRegister = 0x0e
	RegisterMdmcfg1    CxRegister = 0x0f
	RegisterMdmcfg0    CxRegister = 0x10
	RegisterDeviatn    CxRegister = 0x11
	RegisterMcsm0      CxRegister = 0x14
	RegisterFoccfg     CxRegister = 0x15
	RegisterAgcctrl2   CxRegister = 0x17
	RegisterAgcctrl1   CxRegister = 0x18
	RegisterAgcctrl0   CxRegister = 0x19
	RegisterFrend1     CxRegister = 0x1a
	RegisterFrend0     CxRegister = 0x1b
	RegisterFscal3     CxRegister = 0x1c
	RegisterFscal2     CxRegister = 0x1d
	RegisterFscal1     CxRegister = 0x1e
	RegisterFscal0     CxRegister = 0x1f
	RegisterTest1      CxRegister = 0x24
	RegisterTest0      CxRegister = 0x25
	RegisterPaTable0   CxRegister = 0x2e
	// 24MHz crystal
	OscillatorHz = 24000000
)

// GetFrequency returns the radio's current tuning in Hz (from Kenneth)
func (crl *ConnectedRileyLink) GetFrequency() (uint32, error) {
	var (
		frequency uint32 = 0
		value     byte
		err       error
	)

	log.Debug("reading FREQ2")
	value, err = crl.ReadRegister(RegisterFreq2)
	if err != nil {
		return 0, err
	}
	frequency += uint32(value) << 16

	log.Debug("reading FREQ1")
	value, err = crl.ReadRegister(RegisterFreq1)
	if err != nil {
		return 0, err
	}
	frequency += uint32(value) << 8

	log.Debug("reading FREQ0")
	value, err = crl.ReadRegister(RegisterFreq0)
	if err != nil {
		return 0, err
	}
	frequency += uint32(value)

	// oscillator multiplier
	frequency = uint32(uint64(frequency) * OscillatorHz >> 16)

	return frequency, nil
}

// SetFrequency tells the CC to tune to a specific frequency
func (crl *ConnectedRileyLink) SetFrequency(freq uint32) error {
	var err error
	// oscilator multiplier
	freqcal := (uint64(freq)<<16 + OscillatorHz/2) / OscillatorHz

	freq2 := byte(freqcal >> 16)
	log.WithField("freq2", freq2).Debug("writing FREQ2")
	err = crl.WriteRegister(RegisterFreq2, freq2)
	if err != nil {
		return err
	}

	freq1 := byte(freqcal >> 8)
	log.WithField("freq1", freq1).Debug("writing FREQ1")
	err = crl.WriteRegister(RegisterFreq1, freq1)
	if err != nil {
		return err
	}

	freq0 := byte(freqcal)
	log.WithField("freq0", freq0).Debug("writing FREQ0")
	err = crl.WriteRegister(RegisterFreq0, freq0)
	if err != nil {
		return err
	}

	return nil
}
