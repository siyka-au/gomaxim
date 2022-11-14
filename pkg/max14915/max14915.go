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

// Private
func (m *Max14915) read(reg registerAddress, rw commandType) ([]byte, error) {
	return m.transfer(reg, writeCommand, 0x00)
}

func (m *Max14915) write(reg registerAddress, rw commandType, data byte) ([]byte, error) {
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
