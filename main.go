package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/tarm/serial"
)

var serialPort *serial.Port
var serialMu sync.Mutex

func handleConnection(sshSession ssh.Session) {
	serialMu.Lock()
	defer serialMu.Unlock()
	if serialPort == nil {
		fmt.Fprintln(sshSession, "Serial port not initialized")
		return
	}

	// 从SSH会话到串口，同时输出到本地
	go func() {
		multiWriter := io.MultiWriter(serialPort, os.Stdout)
		_, err := io.Copy(multiWriter, sshSession)
		if err != nil {
			log.Println("Error copying from SSH to serial:", err)
		}
	}()

	// 从串口到SSH会话，同时输出到本地
	multiWriter := io.MultiWriter(sshSession, os.Stdout)
	_, err := io.Copy(multiWriter, serialPort)
	if err != nil {
		log.Println("Error copying from serial to SSH:", err)
	}
}

func readSerialPort() {
	c := &serial.Config{Name: "COM3", Baud: 9600} // 还不支持参数，先修改为你的串口号
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

func main() {
	go readSerialPort()

	ssh.Handle(handleConnection)

	log.Println("Starting SSH server on port 2222...")
	err := ssh.ListenAndServe(":2222", nil)
	if err != nil {
		log.Fatal("Failed to start SSH server:", err)
	}
}
