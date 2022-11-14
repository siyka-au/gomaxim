package max14915

import (
	"github.com/siyka-au/gomaxim/internal/maxim"
	"github.com/warthog618/gpiod"
	spidev "golang.org/x/exp/io/spi"
)

type Max14915 struct {
	spi  *spidev.Device
	cs   *gpiod.Line
	addr byte
}

func NewMax14915(spi *spidev.Device, cs *gpiod.Line, addr byte) *Max14915 {
	dev := new(Max14915)
	dev.spi = spi
	dev.cs = cs
	dev.addr = addr
	return dev
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

type Config1 byte

const (
	FaultLEDManualControl    = Config1(0x01)
	StatusLEDManualControl   = Config1(0x02)
	FaultLEDStretchDisable   = Config1(0x00)
	FaultLEDStretch1Sec      = Config1(0x04)
	FaultLEDStretch2Sec      = Config1(0x08)
	FaultLEDStretch3Sec      = Config1(0x0c)
	FaultFilterEnable        = Config1(0x10)
	FilterLongBlankingTime   = Config1(0x20)
	FaultLatchEnable         = Config1(0x40)
	FaultLEDShowCurrentLimit = Config1(0x80)
)

type Config2 byte

const (
	VDDOnThresholdVDDGood                         = Config2(0x01)
	ShortToVDDThreshold9V                         = Config2(0x00)
	ShortToVDDThreshold10V                        = Config2(0x04)
	ShortToVDDThreshold12V                        = Config2(0x08)
	ShortToVDDThreshold14V                        = Config2(0x0c)
	OpenWireCurrentDetectionMagnitude20MicroAmps  = Config2(0x00)
	OpenWireCurrentDetectionMagnitude100MicroAmps = Config2(0x10)
	OpenWireCurrentDetectionMagnitude300MicroAmps = Config2(0x20)
	OpenWireCurrentDetectionMagnitude600MicroAmps = Config2(0x30)
	WatchDogTimeOutDisable                        = Config2(0x00)
	WatchDogTimeOut200Millis                      = Config2(0x40)
	WatchDogTimeOut600Millis                      = Config2(0x80)
	WatchDogTimeOut1_2Secs                        = Config2(0xc0)
)

type Mask byte

const (
	OverloadMask     = Mask(0x01)
	CurrentLimitMask = Mask(0x02)
	OpenWireOffMask  = Mask(0x04)
	OpenWireOnMask   = Mask(0x08)
	ShortToVDDMask   = Mask(0x10)
	VDDOKMask        = Mask(0x20)
	SupplyErrorMask  = Mask(0x40)
	CommErrorMask    = Mask(0x80)
)

func (m *Max14915) ReadOutputs() ([]byte, error) {
	return m.read(outputRegister)
}

func (m *Max14915) WriteOutputs(data Line) ([]byte, error) {
	return m.write(outputRegister, byte(data))
}

func (m *Max14915) ReadGlobalFault() ([]byte, error) {
	return m.read(globalFaultRegister)
}

func (m *Max14915) ReadShortToVDD() ([]byte, error) {
	return m.read(shortToVDDRegister)
}

func (m *Max14915) ReadCurrentLimit() ([]byte, error) {
	return m.read(currentLimitRegister)
}

func (m *Max14915) ReadOverload() ([]byte, error) {
	return m.read(overloadRegister)
}

func (m *Max14915) ReadConfig1() ([]byte, error) {
	return m.read(config1Register)
}

func (m *Max14915) WriteConfig1(data Config1) ([]byte, error) {
	return m.write(config1Register, byte(data))
}

func (m *Max14915) ReadConfig2() ([]byte, error) {
	return m.read(config2Register)
}

func (m *Max14915) WriteConfig2(data Config2) ([]byte, error) {
	return m.write(config2Register, byte(data))
}

func (m *Max14915) ReadMask() ([]byte, error) {
	return m.read(maskRegister)
}

func (m *Max14915) WriteMask(data Mask) ([]byte, error) {
	return m.write(maskRegister, byte(data))
}

// Private
func (m *Max14915) read(reg registerAddress) ([]byte, error) {
	return m.transfer(reg, readCommand, 0x00)
}

func (m *Max14915) write(reg registerAddress, data byte) ([]byte, error) {
	return m.transfer(reg, writeCommand, data)
}

func (m *Max14915) transfer(reg registerAddress, rw commandType, data byte) ([]byte, error) {
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

func (m *Max14915) command(reg registerAddress, rw commandType, burst bool, data byte) ([]byte, error) {
	cmd := make([]byte, 3)
	cmd[0] |= m.addr << 6
	cmd[0] |= byte(reg << 1)
	cmd[0] |= byte(rw)

	cmd[1] = data

	var err error
	cmd[2], err = maxim.CRC(cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

type registerAddress byte

const (
	outputRegister            = registerAddress(0x00)
	faultLEDRegister          = registerAddress(0x01)
	statusLEDRegister         = registerAddress(0x02)
	interruptRegister         = registerAddress(0x03)
	overloadRegister          = registerAddress(0x04)
	currentLimitRegister      = registerAddress(0x05)
	openWireOffRegister       = registerAddress(0x06)
	openWireOnRegister        = registerAddress(0x07)
	shortToVDDRegister        = registerAddress(0x08)
	globalFaultRegister       = registerAddress(0x09)
	openWireOffEnableRegister = registerAddress(0x0a)
	openWireOnEnableRegister  = registerAddress(0x0b)
	shortToVDDEnableRegister  = registerAddress(0x0c)
	config1Register           = registerAddress(0x0d)
	config2Register           = registerAddress(0x0e)
	maskRegister              = registerAddress(0x0f)
)

type commandType byte

const (
	readCommand  = commandType(0)
	writeCommand = commandType(1)
)
