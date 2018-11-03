package gorileylink

type MedtronicPump struct {
	ModelNumber int
}

type MMTPumpSize int

const (
	MMTPumpSizeUnknown MMTPumpSize = 0
	MMTPumpSizeSmall   MMTPumpSize = 500
	MMTPumpSizeLarge   MMTPumpSize = 700
)

func (mmtpump *MedtronicPump) GetGeneration() int {
	return mmtpump.ModelNumber
}

func (mmtpump *MedtronicPump) GetMaxReserviorSize() MMTPumpSize {
	if mmtpump.ModelNumber&int(MMTPumpSizeSmall) == int(MMTPumpSizeSmall) {
		return MMTPumpSizeSmall
	} else if mmtpump.ModelNumber&int(MMTPumpSizeLarge) == int(MMTPumpSizeLarge) {
		return MMTPumpSizeLarge
	} else {
		return MMTPumpSizeUnknown
	}
}

func (mmtpump *MedtronicPump) generation() int {
	return mmtpump.ModelNumber % 100
}

// feature flags of the pump

func (mmtpump *MedtronicPump) NewRecordStyle() bool {
	return mmtpump.Modern()
}

func (mmtpump *MedtronicPump) ASWTHOSOD() bool {
	// appendsSquareWaveToHistoryOnStartOfDelivery
	return mmtpump.Modern()
}

func (mmtpump *MedtronicPump) HasMySentry() bool {
	return mmtpump.Modern()
}

func (mmtpump *MedtronicPump) HasLowSuspend() bool {
	return mmtpump.generation() >= 51
}

func (mmtpump *MedtronicPump) RBPSE() bool {
	// recordsBasalProfileStartEvents
	return mmtpump.Modern()
}

func (mmtpump *MedtronicPump) HasBolusErrorQuirk() bool {
	// // On x15 models, a bolus in progress error is returned when bolusing, even though the bolus succeeds
	return mmtpump.Modern()
}

// Modern returns whether the pump is considered "modern"
func (mmtpump *MedtronicPump) Modern() bool {
	return mmtpump.generation() >= 23
}

func (mmtpump *MedtronicPump) StrokesPerUnit() int {
	// /// Newer models allow higher precision delivery, and have bit packing to accomodate this.
	if mmtpump.Modern() {
		return 40
	}
	return 10
}

var knownPumps = map[int]string{
	508: "508",
	511: "511",
	711: "711",
	512: "512",
	712: "712",
	515: "515",
	715: "715",
	522: "522",
	722: "722",
	523: "523",
	723: "723",
	530: "530",
	730: "730",
	540: "540",
	740: "740",
	551: "551",
	751: "751",
	554: "554",
	754: "754",
}
