package gorileylink

type LEDMode byte

const (
	LEDModeOff  LEDMode = 0x00
	LEDModeOn   LEDMode = 0x01
	LEDModeAuto LEDMode = 0x02
)
