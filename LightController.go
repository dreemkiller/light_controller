package main

import (
	"fmt"
	"time"
    "github.com/stianeikeland/go-rpio"
)

func main() {
	fmt.Printf("Hello, world\n")
	var programs [1]Program
	if err := programs[0].Load("program0"); err != nil {
		fmt.Printf("Failed to load:%v\n", err)
		return
	}
	fmt.Printf("programs[0]:%v\n", programs[0])

	if err := rpio.Open(); err != nil {
		fmt.Printf("rpio.Open() failed:%v\n", err)
		return
	}
	var pins [8]rpio.Pin
	pins[0] = rpio.Pin(2)
	pins[1] = rpio.Pin(3)
	pins[2] = rpio.Pin(4)
	pins[3] = rpio.Pin(17)
	pins[4] = rpio.Pin(27)
	pins[5] = rpio.Pin(22)
	pins[6] = rpio.Pin(10)
	pins[7] = rpio.Pin(9)
	for _, pin := range pins {
		pin.Output()
	}

	for {
		pause := time.Millisecond * time.Duration(programs[0].timeslice_ms)
		for _, this_pin := range pins {
			this_pin.Low()
		}
		for moment_count, moment := range programs[0].moments {
			fmt.Printf("slice:%v\n", moment_count)
			for pin_num, value := range moment.lights {
				if (value) {
					pins[pin_num].High()
				} else {
					pins[pin_num].Low()
				}
			}
			time.Sleep(pause)
		}
	}
}
