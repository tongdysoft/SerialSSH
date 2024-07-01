package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

func readSerialPort() {
	c := &serial.Config{Name: "COM1", Baud: 9600} // 还不支持参数，先修改为你的串口号
	log.Printf("%s: %s, %s: %d, %s: %d, %s: %s, %s: %d, %s: %d\n", l("SERIALPORT"), c.Name, l("BAUD"), c.Baud, l("SIZE"), c.Size, l("PARITY"), string(c.Parity), l("STOPBITS"), c.StopBits, l("TIMEOUT"), c.ReadTimeout)
	var err error
	serialPort, err = serial.OpenPort(c)
	if err != nil {
		log.Fatalf("%s: %v", l("OPENSERIALERR"), err)
	}
	defer serialPort.Close()

	buf := make([]byte, 128)
	for {
		n, err := serialPort.Read(buf)
		if err != nil {
			log.Printf("%s: %v", l("SERIALREADERR"), err)
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Println(string(buf[:n]))
	}
}
