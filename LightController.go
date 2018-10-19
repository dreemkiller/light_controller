package main

import (
	"fmt"
	"time"
    "github.com/stianeikeland/go-rpio"
)

func main() {
	fmt.Printf("Hello, world\n")
	var program Program
	err := program.Load("goober"); if err != nil {
		fmt.Printf("Failed to load:%v\n", err)
		return
	}
	fmt.Printf("program:%v\n", program)

	pin := rpio.Pin(4)
	pin.Output()
	i := 0
	pause := time.Second
	for  i < 100 {
		i += 1
		time.Sleep(pause)
		fmt.Printf("On\n")
		pin.High()
		time.Sleep(pause)
		fmt.Printf("Off\n")
		pin.Low()
	}
}
