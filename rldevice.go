package gorileylink

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

// RileyLinkCommand is the literal type of a device command
type RileyLinkCommand byte

const (
	RLCGetState         RileyLinkCommand = 0x01
	RLCGetVersion       RileyLinkCommand = 0x02
	RLCGetPacket        RileyLinkCommand = 0x03
	RLCSendPacket       RileyLinkCommand = 0x04
	RLCSendAndListen    RileyLinkCommand = 0x05
	RLCUpdateRegister   RileyLinkCommand = 0x06
	RLCReset            RileyLinkCommand = 0x07
	RLCLED              RileyLinkCommand = 0x08
	RLCReadRegister     RileyLinkCommand = 0x09
	RLCSetModeRegisters RileyLinkCommand = 0x10
	RLCSetSWEncoding    RileyLinkCommand = 0x11
	RLCSetPreamble      RileyLinkCommand = 0x12
	RLCResetRadioConfig RileyLinkCommand = 0x13
	RLCGetStatistics    RileyLinkCommand = 0x14
)
