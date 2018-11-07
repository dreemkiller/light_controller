package main

import "os"
import "fmt"
import "io"
import "time"
import "github.com/stianeikeland/go-rpio"

type Program struct {
	moments []moment
	timeslice_ms int
}

type moment struct {
	lights [8]bool
}

func (p *Program) Load(filename string) error {
	fd, err := os.Open(filename); if err != nil {
		return err
	}
	if _, err := fmt.Fscanf(fd, "%d\n", &p.timeslice_ms); err != nil {
		fmt.Printf("Failed to read timeslice from program file:%v\n", err)
		return err
	}
	for {
		var lights [8]int
		_, err := fmt.Fscanf(fd, "%1d%1d%1d%1d%1d%1d%1d%1d\n", &lights[0],
															   &lights[1],
															   &lights[2],
															   &lights[3],
															   &lights[4],
															   &lights[5],
															   &lights[6],
															   &lights[7])
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		var this_moment moment
		for i, light := range lights {
			this_moment.lights[i] = (light == 1)
		}
		p.moments = append(p.moments, this_moment)
	}
	return nil
}

func (p *Program) Run(pins []rpio.Pin) {
	pause := time.Millisecond * time.Duration(p.timeslice_ms)
	for _, this_pin := range pins {
		this_pin.High()
	}
	for moment_count, moment := range p.moments {
		fmt.Printf("slice:%v\n", moment_count)
		for pin_num, value := range moment.lights {
			if (value) {
				pins[pin_num].Low()
			} else {
				pins[pin_num].High()
			}
		}
		time.Sleep(pause)
	}
}
