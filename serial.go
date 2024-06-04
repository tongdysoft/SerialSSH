package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

func readSerialPort() {
	c := &serial.Config{Name: "COM1", Baud: 9600} // 还不支持参数，先修改为你的串口号
	var err error
	serialPort, err = serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	defer serialPort.Close()

	buf := make([]byte, 128)
	for {
		n, err := serialPort.Read(buf)
		if err != nil {
			log.Printf("Error reading from serial port: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Println(string(buf[:n]))
	}
}
