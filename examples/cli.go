package main

import (
	"flag"
	"log"
	"os"

	gpio "github.com/DoZator/go_gpio"
)

var (
	number = flag.Int("n", 2, "GPIO pin number")
	input  = flag.Bool("in", false, "Set input pin mode")
	out    = flag.Bool("out", true, "Set output pin mode")
	read   = flag.Bool("read", false, "Get current pin state")
	high   = flag.Bool("high", false, "Set pin high")
	low    = flag.Bool("low", true, "Set pin low")
	clean  = flag.Bool("clean", false, "Cleanup pins")
)

func main() {
	flag.Parse()

	mode := gpio.ModeOUT
	if *input {
		mode = gpio.ModeIN
	}

	pin := gpio.Setup(*number, mode)

	if *clean {
		gpio.Cleanup()
		os.Exit(0)
	}

	if *read {
		v := pin.Read()
		log.Println(v)
		os.Exit(0)
	}

	if *high {
		gpio.Output(*number, gpio.PinHIGH)
		os.Exit(0)
	}

	if *low {
		pin.Low()
	}
}
