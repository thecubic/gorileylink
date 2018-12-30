package gorileylink

import (
	"time"
)

// LEDMode is the literal type of the LED mode flag
type LEDMode byte

const (
	// LEDModeOff turns diag LEDs off
	LEDModeOff LEDMode = 0x00
	// LEDModeOn turns diag LEDs on
	LEDModeOn LEDMode = 0x01
	// LEDModeAuto ???
	LEDModeAuto LEDMode = 0x02
)

type LEDColor byte

const (
	LEDGreen LEDColor = 0x00
	LEDBlue  LEDColor = 0x01
)

// RileyLinkCommand is the literal type of a device command
type RileyLinkCommand byte

const (
	RLCInterrupt        RileyLinkCommand = 0x00
	RLCGetState         RileyLinkCommand = 0x01
	RLCGetVersion       RileyLinkCommand = 0x02
	RLCGetPacket        RileyLinkCommand = 0x03
	RLCSendPacket       RileyLinkCommand = 0x04
	RLCSendAndListen    RileyLinkCommand = 0x05
	RLCUpdateRegister   RileyLinkCommand = 0x06
	RLCReset            RileyLinkCommand = 0x07
	RLCLED              RileyLinkCommand = 0x08
	RLCReadRegister     RileyLinkCommand = 0x09
	RLCSetModeRegisters RileyLinkCommand = 0x0A
	RLCSetSWEncoding    RileyLinkCommand = 0x0B
	RLCSetPreamble      RileyLinkCommand = 0x0C
	RLCResetRadioConfig RileyLinkCommand = 0x0D
	RLCGetStatistics    RileyLinkCommand = 0x0E
)

// RileyLinkCCResponseType represents the outcome of the sent command
type RileyLinkCCResponseType byte

const (
	RLRRecvTimeout    RileyLinkCCResponseType = 0xaa
	RLRInterrupted    RileyLinkCCResponseType = 0xbb
	RLRZeroData       RileyLinkCCResponseType = 0xcc
	RLRSuccess        RileyLinkCCResponseType = 0xdd
	RLRInvalidParam   RileyLinkCCResponseType = 0x11
	RLRUnknownCommand RileyLinkCCResponseType = 0x22
)

func (rlr RileyLinkCCResponseType) String() string {
	if rlr == RLRRecvTimeout {
		return "RLRRecvTimeout"
	} else if rlr == RLRInterrupted {
		return "RLRInterrupted"
	} else if rlr == RLRZeroData {
		return "RLRZeroData"
	} else if rlr == RLRSuccess {
		return "RLRSuccess"
	} else if rlr == RLRInvalidParam {
		return "RLRInvalidParam"
	} else if rlr == RLRUnknownCommand {
		return "RLRUnknownCommand"
	}
	return "RLRUnknown"
}

type RileyLinkPacketChannel byte

const (
	RLPCCGM  RileyLinkPacketChannel = 0x01
	RLPCPUMP RileyLinkPacketChannel = 0x02
)

// RileyLinkStatistics represents a statistics pull event
type RileyLinkStatistics struct {
	Uptime            time.Duration
	RecvOverflows     uint16
	RecvFifoOverflows uint16
	PacketsRecv       uint16
	PacketsXmit       uint16
	CRCFailures       uint16
	SPISyncFailures   uint16
}
