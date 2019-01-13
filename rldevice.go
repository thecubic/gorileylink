package gorileylink

import (
	"time"
)

// LEDMode is the literal type of the LED mode flag
type LEDMode byte

const (
	// LEDOff turns diag LEDs off
	LEDOff LEDMode = 0x00
	// LEDOn turns diag LEDs on
	LEDOn LEDMode = 0x01
	// LEDAuto ???
	LEDAuto LEDMode = 0x02
)

func (ledm LEDMode) String() string {
	switch ledm {
	case LEDOff:
		return "LEDOff"
	case LEDOn:
		return "LEDOn"
	case LEDAuto:
		return "LEDAuto"
	default:
		return "LEDModeUNKNOWN"
	}
}

// LEDColor is the literal type of the choice of LED
// Valid for CC LEDs for sure, BLE LEDs idk
type LEDColor byte

const (
	// LEDGreen is Green, really zeroth
	LEDGreen LEDColor = 0x00
	// LEDBlue is Blue, really first
	LEDBlue LEDColor = 0x01
)

func (ledc LEDColor) String() string {
	switch ledc {
	case LEDGreen:
		return "LEDGreen"
	case LEDBlue:
		return "LEDBlue"
	default:
		return "LedColorUNKNOWN"
	}
}

// RileyLinkCommand is the literal type of a device command
type RileyLinkCommand byte

const (
	// RLCInterrupt just pushes a null command out, effectively
	// interrupting another command in progress
	RLCInterrupt RileyLinkCommand = 0x00
	// RLCGetState returns "OK" when the CC chip is ... OK
	RLCGetState RileyLinkCommand = 0x01
	// RLCGetVersion returns the version of the CC firmware (subg_rfspy)
	RLCGetVersion RileyLinkCommand = 0x02

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

func (rlc RileyLinkCommand) String() string {
	switch rlc {
	case RLCInterrupt:
		return "RLCInterrupt"
	case RLCGetState:
		return "RLCGetState"
	case RLCGetVersion:
		return "RLCGetVersion"
	case RLCGetPacket:
		return "RLCGetPacket"
	case RLCSendPacket:
		return "RLCSendPacket"
	case RLCSendAndListen:
		return "RLCSendAndListen"
	case RLCUpdateRegister:
		return "RLCUpdateRegister"
	case RLCReset:
		return "RLCReset"
	case RLCLED:
		return "RLCLED"
	case RLCReadRegister:
		return "RLCReadRegister"
	case RLCSetModeRegisters:
		return "RLCSetModeRegisters"
	case RLCSetSWEncoding:
		return "RLCSetSWEncoding"
	case RLCSetPreamble:
		return "RLCSetPreamble"
	case RLCResetRadioConfig:
		return "RLCResetRadioConfig"
	case RLCGetStatistics:
		return "RLCGetStatistics"
	default:
		return "RileyLinkCommandUNKNOWN"
	}
}

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
	switch rlr {
	case RLRRecvTimeout:
		return "RLRRecvTimeout"
	case RLRInterrupted:
		return "RLRInterrupted"
	case RLRZeroData:
		return "RLRZeroData"
	case RLRSuccess:
		return "RLRSuccess"
	case RLRInvalidParam:
		return "RLRInvalidParam"
	case RLRUnknownCommand:
		return "RLRUnknownCommand"
	default:
		return "RileyLinkCCResponseTypeUNKNOWN"
	}
}

// RileyLinkPacketChannel represents the internal channel type
type RileyLinkPacketChannel byte

const (
	RLPCMeter RileyLinkPacketChannel = 0x01
	RLPCPump  RileyLinkPacketChannel = 0x02
)

func (rlpc RileyLinkPacketChannel) String() string {
	switch rlpc {
	case RLPCMeter:
		return "RLPCMeter"
	case RLPCPump:
		return "RLPCPump"
	default:
		return "RileyLinkPacketChannelUNKNOWN"
	}
}

// RileyLinkStatistics represents a concrete statistics pull event
type RileyLinkStatistics struct {
	Collected         time.Time
	Uptime            time.Duration
	RecvOverflows     uint16
	RecvFifoOverflows uint16
	PacketsRecv       uint16
	PacketsXmit       uint16
	CRCFailures       uint16
	SPISyncFailures   uint16
}
