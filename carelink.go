package gorileylink

// Carelink describes the Medtronic exchange protocol

type CarelinkMessageType byte

type CarelinkMessage struct {
	MessageType CarelinkMessageType
	Data        []byte
}

const (
	CMTAlert                        CarelinkMessageType = 0x01
	CMTAlertCleared                 CarelinkMessageType = 0x02
	CMTDeviceTest                   CarelinkMessageType = 0x03
	CMTPumpStatus                   CarelinkMessageType = 0x04
	CMTPumpAck                      CarelinkMessageType = 0x06
	CMTPumpBackfill                 CarelinkMessageType = 0x08
	CMTFindDevice                   CarelinkMessageType = 0x09
	CMTDeviceLink                   CarelinkMessageType = 0x0A
	CMTErrorResponse                CarelinkMessageType = 0x15
	CMTWriteGlucoseHistoryTimestamp CarelinkMessageType = 0x28

	CMTSetBasalProfileA CarelinkMessageType = 0x30 // CMD_SET_A_PROFILE
	CMTSetBasalProfileB CarelinkMessageType = 0x31 // CMD_SET_B_PROFILE

	CMTChangeTime  CarelinkMessageType = 0x40
	CMTSetMaxBolus CarelinkMessageType = 0x41 // CMD_SET_MAX_BOLUS
	CMTBolus       CarelinkMessageType = 0x42

	CMTPumpExperiment_OP67 CarelinkMessageType = 0x43
	CMTPumpExperiment_OP68 CarelinkMessageType = 0x44
	CMTPumpExperiment_OP69 CarelinkMessageType = 0x45 // CMD_SET_VAR_BOLUS_ENABLE

	CMTSelectBasalProfile CarelinkMessageType = 0x4a

	CMTChangeTempBasal CarelinkMessageType = 0x4c

	CMTPumpExperiment_OP80     CarelinkMessageType = 0x50
	CMTSetRemoteControlID      CarelinkMessageType = 0x51 // CMD_SET_RF_REMOTE_ID
	CMTPumpExperiment_OP82     CarelinkMessageType = 0x52 // CMD_SET_BLOCK_ENABLE
	CMTSetLanguage             CarelinkMessageType = 0x53
	CMTPumpExperiment_OP84     CarelinkMessageType = 0x54 // CMD_SET_ALERT_TYPE
	CMTPumpExperiment_OP85     CarelinkMessageType = 0x55 // CMD_SET_PATTERNS_ENABLE
	CMTPumpExperiment_OP86     CarelinkMessageType = 0x56
	CMTSetRemoteControlEnabled CarelinkMessageType = 0x57 // CMD_SET_RF_ENABLE
	CMTPumpExperiment_OP88     CarelinkMessageType = 0x58 // CMD_SET_INSULIN_ACTION_TYPE
	CMTPumpExperiment_OP89     CarelinkMessageType = 0x59
	CMTPumpExperiment_OP90     CarelinkMessageType = 0x5a

	CMTButtonPress CarelinkMessageType = 0x5b

	CMTPumpExperiment_OP92 CarelinkMessageType = 0x5c

	CMTPowerOn CarelinkMessageType = 0x5d

	CMTSetBolusWizardEnabled1 CarelinkMessageType = 0x61
	CMTSetBolusWizardEnabled2 CarelinkMessageType = 0x62
	CMTSetBolusWizardEnabled3 CarelinkMessageType = 0x63
	CMTSetBolusWizardEnabled4 CarelinkMessageType = 0x64
	CMTSetBolusWizardEnabled5 CarelinkMessageType = 0x65
	CMTSetAlarmClockEnable    CarelinkMessageType = 0x67

	CMTSetMaxBasalRate         CarelinkMessageType = 0x6e // CMD_SET_MAX_BASAL
	CMTSetBasalProfileStandard CarelinkMessageType = 0x6f // CMD_SET_STD_PROFILE

	CMTReadTime             CarelinkMessageType = 0x70
	CMTGetBattery           CarelinkMessageType = 0x72
	CMTReadRemainingInsulin CarelinkMessageType = 0x73
	CMTReadFirmwareVersion  CarelinkMessageType = 0x74
	CMTReadErrorStatus      CarelinkMessageType = 0x75
	CMTReadRemoteControlIDs CarelinkMessageType = 0x76 // CMD_READ_REMOTE_CTRL_IDS

	CMTGetHistoryPage         CarelinkMessageType = 0x80
	CMTGetPumpModel           CarelinkMessageType = 0x8d
	CMTReadProfileSTD512      CarelinkMessageType = 0x92
	CMTReadProfileA512        CarelinkMessageType = 0x93
	CMTReadProfileB512        CarelinkMessageType = 0x94
	CMTReadTempBasal          CarelinkMessageType = 0x98
	CMTGetGlucosePage         CarelinkMessageType = 0x9A
	CMTReadCurrentPageNumber  CarelinkMessageType = 0x9d
	CMTReadSettings           CarelinkMessageType = 0xc0
	CMTReadCurrentGlucosePage CarelinkMessageType = 0xcd
	CMTReadPumpStatus         CarelinkMessageType = 0xce

	CMTUnknown_e2            CarelinkMessageType = 0xe2 // a7594040e214190226330000000000021f99011801e00103012c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
	CMTUnknown_e6            CarelinkMessageType = 0xe6 // a7594040e60200190000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
	CMTSettingsChangeCounter CarelinkMessageType = 0xec // Body[3] increments by 1 after changing certain settings 0200af0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000

	CMTReadOtherDevicesIDs      CarelinkMessageType = 0xf0
	CMTReadCaptureEventEnabled  CarelinkMessageType = 0xf1 // Body[1] encodes the bool state 0101000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
	CMTChangeCaptureEventEnable CarelinkMessageType = 0xf2
	CMTReadOtherDevicesStatus   CarelinkMessageType = 0xf3
)
