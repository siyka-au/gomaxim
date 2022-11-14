package main

import (
	spidev "golang.org/x/exp/io/spi"
)

func main() {
	// GPIO setup

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

	// Rah
}
