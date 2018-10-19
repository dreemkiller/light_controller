package main

import "os"
import "fmt"
import "io"

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
	fmt.Printf("p:%v\n", p)
	return nil
}
