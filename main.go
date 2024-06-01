package main

import (
	"fmt"
	"io"
	"log"

	"github.com/gliderlabs/ssh"
	"github.com/tarm/serial"
)

func handleConnection(sshSession ssh.Session) {
	// 配置串口
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Fprintln(sshSession, "Failed to open serial port:", err)
		return
	}
	defer s.Close()

	// 创建双向管道
	go func() {
		_, err := io.Copy(s, sshSession) // 从SSH会话到串口
		if err != nil {
			log.Println("Error copying from SSH to serial:", err)
		}
	}()
	_, err = io.Copy(sshSession, s) // 从串口到SSH会话
	if err != nil {
		log.Println("Error copying from serial to SSH:", err)
	}
}

func main() {
	ssh.Handle(handleConnection)

	log.Println("Starting SSH server on port 2222...")
	err := ssh.ListenAndServe(":2222", nil)
	if err != nil {
		log.Fatal("Failed to start SSH server:", err)
	}
}
