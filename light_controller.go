package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"os"
	"io/ioutil"
    "github.com/stianeikeland/go-rpio"
)

type Config struct {
	Port int	`json:"port"`
	Server string	`json:"server"`
	Cert_file string `json:"cert_file"`
}

func readConfig() error {
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Printf("Failed to open config file:%v\n", err)
		return err
	}
	defer jsonFile.Close()
	jsonParser := json.NewDecoder(jsonFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		fmt.Printf("Failed to decode json:%v\n", err)
		return err
	}
	return nil
}

var current_program_num int

var mutex *sync.Mutex
var programs [3]Program

var config Config

func main() {
	fmt.Printf("Hello, world\n")
	err := readConfig()	
	if err != nil {
		fmt.Printf("readConfig failed:%v\n", err)
		return
	}
	fmt.Printf("config.Server:%v\n", config.Server)
	fmt.Printf("config.Port:%v\n", config.Port)
	mutex = &sync.Mutex{}
	if err := programs[0].Load("program0"); err != nil {
		fmt.Printf("Failed to load:%v\n", err)
		return
	}
	if err := programs[1].Load("program1"); err != nil {
		fmt.Printf("Failed to load:%v\n", err)
		return
	}
	if err := programs[2].Load("program2"); err != nil {
		fmt.Printf("Failed to load program 2:%v\n", err)
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
		var temp_program_num int
		mutex.Lock()
		temp_program_num = current_program_num
		mutex.Unlock()
		programs[temp_program_num].Run(pins[:])
	}
}

type ProgramNumber struct {
	Number int `json:"Number"`
}

const fetch_delay = 5

func GetProgram() {
	server_url := fmt.Sprintf("https://%s:%d/CurrentProgram", config.Server, config.Port)
	fmt.Printf("server_url:%v\n", server_url)
	caCert, err := ioutil.ReadFile(config.Cert_file)
	if err != nil {
		fmt.Printf("Failed to read server crt:%v\n", err)
		return		
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client := &http.Client {
		Transport: &http.Transport {
			TLSClientConfig: &tls.Config {
				RootCAs: caCertPool,
			},
		},
	}
	for {
		var received_program_num ProgramNumber
		time.Sleep(time.Second * time.Duration(fetch_delay))
		resp, err := client.Get(server_url)
		if err != nil {
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
