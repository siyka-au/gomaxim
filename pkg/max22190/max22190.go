package max22190

import (
	"github.com/siyka-au/gomaxim/internal/maxim"
	"github.com/warthog618/gpiod"
	spidev "golang.org/x/exp/io/spi"
)

type Max22190 struct {
	spi *spidev.Device
	cs  *gpiod.Line
}

// Public
type Line byte

const (
	Line1 = Line(0x01)
	Line2 = Line(0x02)
	Line3 = Line(0x04)
	Line4 = Line(0x08)
	Line5 = Line(0x10)
	Line6 = Line(0x20)
	Line7 = Line(0x40)
	Line8 = Line(0x80)
)

type Fault1 byte

const (
	WireBreakFault             = Fault1(0x01)
	VoltageThreshold24VMFault  = Fault1(0x02)
	VoltageThreshold24VLFault  = Fault1(0x04)
	TemperatureThreshold1Alarm = Fault1(0x08)
	TemperatureThreshold2Alarm = Fault1(0x10)
	Fault2Active               = Fault1(0x20)
	PowerOnReset               = Fault1(0x40)
	CRCFault                   = Fault1(0x80)
)

type Fault1Enable byte

const (
	WireBreakFaultEnable             = Fault1Enable(0x01)
	VoltageThreshold24VMFaultEnable  = Fault1Enable(0x02)
	VoltageThreshold24VLFaultEnable  = Fault1Enable(0x04)
	TemperatureThreshold1AlarmEnable = Fault1Enable(0x08)
	TemperatureThreshold2AlarmEnable = Fault1Enable(0x10)
	Fault2ActiveEnable               = Fault1Enable(0x20)
	PowerOnResetEnable               = Fault1Enable(0x40)
	CRCFaultEnable                   = Fault1Enable(0x80)
)

type Filter byte

const (
	Delay50Micros            = Filter(0x00)
	Delay100Micros           = Filter(0x01)
	Delay400Micros           = Filter(0x02)
	Delay800Micros           = Filter(0x03)
	Delay1_6Millis           = Filter(0x04)
	Delay3_2Millis           = Filter(0x05)
	Delay12_8Millis          = Filter(0x06)
	Delay20_6Millis          = Filter(0x07)
	ProgrammableFilterBypass = Filter(0x08)
	WireBreakDetectionEnable = Filter(0x10)
)

type Configuration byte

const (
	REFDIShortCircuitDetectionEnable           = Configuration(0x01)
	AllFiltersFixedToMidScale                  = Configuration(0x08)
	VoltageThresholdFaultsRequireExplicitClear = Configuration(0x10)
)

type Fault2 byte

const (
	REFWBShortCircuitFault  = Fault2(0x01)
	REFWBOpenCircuitFault   = Fault2(0x02)
	REFDIShortCircuitFault  = Fault2(0x04)
	REFDIOpenCircuitFault   = Fault2(0x08)
	OverTemperatureShutdown = Fault2(0x10)
	ClockNot8PulsesFault    = Fault2(0x20)
)

type Fault2Enable byte

const (
	REFWBShortCircuitFaultEnable  = Fault2Enable(0x01)
	REFWBOpenCircuitFaultEnable   = Fault2Enable(0x02)
	REFDIShortCircuitFaultEnable  = Fault2Enable(0x04)
	REFDIOpenCircuitFaultEnable   = Fault2Enable(0x08)
	OverTemperatureShutdownEnable = Fault2Enable(0x10)
	ClockNot8PulsesFaultEnable    = Fault2Enable(0x20)
)

type GPO byte

const (
	FaultStickyPin = GPO(0x80)
)

func (m *Max22190) WriteFault1(reg registerAddress, rw commandType, data byte) ([]byte, error) {
	m.cs.SetValue(0)
	cmd, err := m.command(reg, rw, false, data)
	if err != nil {
		return nil, err
	}
	resp := make([]byte, 3)
	if err := m.spi.Tx(cmd, resp); err != nil {
		return nil, err
	}
	m.cs.SetValue(1)
	return resp, nil
}

// Private
type registerAddress byte

const (
	wireBreakRegister     = registerAddress(0x00)
	digitalInputRegister  = registerAddress(0x02)
	fault1Register        = registerAddress(0x04)
	filterIn1Register     = registerAddress(0x06)
	filterIn2Register     = registerAddress(0x08)
	filterIn3Register     = registerAddress(0x0a)
	filterIn4Register     = registerAddress(0x0c)
	filterIn5Register     = registerAddress(0x0e)
	filterIn6Register     = registerAddress(0x10)
	filterIn7Register     = registerAddress(0x12)
	filterIn8Register     = registerAddress(0x14)
	configurationRegister = registerAddress(0x18)
	inputEnableRegister   = registerAddress(0x1a)
	fault2Register        = registerAddress(0x1c)
	fault2EnableRegister  = registerAddress(0x1e)
	gpoRegister           = registerAddress(0x22)
	fault1EnableRegister  = registerAddress(0x24)
	noOpRegister          = registerAddress(0x26)
)

type commandType byte

const (
	readCommand  = commandType(0)
	writeCommand = commandType(1)
)

func (m *Max22190) read(reg registerAddress, rw commandType) ([]byte, error) {
	return m.transfer(reg, writeCommand, 0x00)
}

func (m *Max22190) write(reg registerAddress, rw commandType, data byte) ([]byte, error) {
	return m.transfer(reg, writeCommand, data)
}

func (m *Max22190) transfer(reg registerAddress, rw commandType, data byte) ([]byte, error) {
	m.cs.SetValue(0)
	cmd, err := m.command(reg, rw, false, data)
	if err != nil {
		return nil, err
	}
	resp := make([]byte, 3)
	if err := m.spi.Tx(cmd, resp); err != nil {
		return nil, err
	}
	m.cs.SetValue(1)
	return resp, nil
}

func (m *Max22190) command(reg registerAddress, rw commandType, burst bool, data byte) ([]byte, error) {
	cmd := make([]byte, 3)
	cmd[0] |= byte(reg)
	cmd[0] |= byte(rw << 7)

	cmd[1] = data

	var err error
	cmd[2], err = maxim.CRC(cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
