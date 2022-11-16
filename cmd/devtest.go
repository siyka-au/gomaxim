package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/siyka-au/gomaxim/pkg/max14915"
	"github.com/siyka-au/gomaxim/pkg/max22190"
	"github.com/warthog618/gpiod"
	spidev "golang.org/x/exp/io/spi"
)

func main() {
	// GPIO setup
	fault, err := gpiod.RequestLine("gpiochip1", 10, gpiod.AsInput)
	if err != nil {
		panic(err)
	}

	latch, err := gpiod.RequestLine("gpiochip2", 10, gpiod.AsOutput(1))
	if err != nil {
		panic(err)
	}

	sync, err := gpiod.RequestLine("gpiochip2", 9, gpiod.AsOutput(1))
	if err != nil {
		panic(err)
	}

	cs1, err := gpiod.RequestLine("gpiochip8", 3, gpiod.AsOutput(1))
	if err != nil {
		panic(err)
	}

	cs2, err := gpiod.RequestLine("gpiochip2", 12, gpiod.AsOutput(1))
	if err != nil {
		panic(err)
	}

	cs3, err := gpiod.RequestLine("gpiochip1", 12, gpiod.AsOutput(1))
	if err != nil {
		panic(err)
	}

	cs4, err := gpiod.RequestLine("gpiochip2", 11, gpiod.AsOutput(1))
	if err != nil {
		panic(err)
	}

	// revert line to input on the way out.
	defer func() {
		fault.Close()
		latch.Reconfigure(gpiod.AsInput)
		latch.Close()
		sync.Reconfigure(gpiod.AsInput)
		sync.Close()
		cs1.Reconfigure(gpiod.AsInput)
		cs1.Close()
		cs2.Reconfigure(gpiod.AsInput)
		cs2.Close()
		cs3.Reconfigure(gpiod.AsInput)
		cs3.Close()
		cs4.Reconfigure(gpiod.AsInput)
		cs4.Close()
	}()

	// SPI setup
	spi, err := spidev.Open(&spidev.Devfs{
		Dev:      "/dev/spidev0.0",
		Mode:     spidev.Mode0,
		MaxSpeed: 1_000_000,
	})
	if err != nil {
		panic(err)
	}
	defer spi.Close()

	// Device initialisation
	in1 := max22190.NewMax22190(spi, cs1)
	in2 := max22190.NewMax22190(spi, cs2)
	in3 := max22190.NewMax22190(spi, cs3)
	out1 := max14915.NewMax14915(spi, cs4, 0b00)
	out2 := max14915.NewMax14915(spi, cs4, 0b01)

	// Capture exit signals to ensure pin is reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	faults, _ := in1.ReadFault1()
	fmt.Printf("IN1 Faults: %08b\n", faults[1])

	faults, _ = in2.ReadFault1()
	fmt.Printf("IN2 Faults: %08b\n", faults[1])

	faults, _ = in3.ReadFault1()
	fmt.Printf("IN3 Faults: %08b\n", faults[1])

	faults, _ = out1.ReadGlobalFault()
	fmt.Printf("OUT1 Global Fault: %08b\n", faults[1])
	conf, _ := out1.ReadConfig1()
	fmt.Printf("OUT1 Config1: %08b\n", conf[1])
	conf, _ = out1.ReadConfig2()
	fmt.Printf("OUT1 Config2: %08b\n", conf[1])

	faults, _ = out1.ReadGlobalFault()
	fmt.Printf("OUT2 Global Fault: %b\n", faults[1])
	conf, _ = out2.ReadConfig1()
	fmt.Printf("OUT2 Config1: %08b\n", conf[1])
	conf, _ = out2.ReadConfig2()
	fmt.Printf("OUT2 Config2: %08b\n", conf[1])

	out1.WriteConfig1(max14915.FaultLatchEnable | max14915.FaultFilterEnable)
	out2.WriteConfig1(max14915.FaultLatchEnable | max14915.FaultFilterEnable)

	v := 0
	for {
		select {
		case <-time.After(500 * time.Millisecond):
			v ^= 1
			if v == 1 {
				out1.SetOutputs(max14915.Line1 | max14915.Line4)
			} else {
				out1.ResetOutputs(max14915.Line1 | max14915.Line4)
			}
		case <-quit:
			return
		}
	}

}
