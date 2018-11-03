package gorileylink

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
)
