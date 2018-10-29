package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
    "github.com/stianeikeland/go-rpio"
)

const server_url = "http://192.168.0.110:8000/CurrentProgram"

var current_program_num int

var mutex *sync.Mutex
var programs [2]Program

func main() {
	fmt.Printf("Hello, world\n")
	mutex = &sync.Mutex{}
	if err := programs[0].Load("program0"); err != nil {
		fmt.Printf("Failed to load:%v\n", err)
		return
	}
	if err := programs[1].Load("program1"); err != nil {
		fmt.Printf("Failed to load:%v\n", err)
		return
	}

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

	go GetProgram()

	for {
		programs[current_program_num].Run(pins[:])
	}
}

type ProgramNumber struct {
	Number int `json:"Number"`
}

func GetProgram() {
	for {
		var received_program_num ProgramNumber
		time.Sleep(time.Second * time.Duration(5))
		resp, err := http.Get(server_url); if err != nil {
			fmt.Printf("Failed to get:%v\n", err)
			continue
		}
		json.NewDecoder(resp.Body).Decode(&received_program_num)
		resp.Body.Close()
		fmt.Printf("received_program_num:%v\n", received_program_num.Number)
		if (received_program_num.Number != current_program_num &&
		    received_program_num.Number < len(programs)) {
			mutex.Lock()
			current_program_num = received_program_num.Number
			mutex.Unlock()
		}
	}
}
