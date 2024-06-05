package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

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

func main() {
	go readSerialPort()

	privateKey, err := loadOrGenerateSSHKey("server.pem")
	if err != nil {
		log.Fatalf("Failed to load or generate private key: %v", err)
	}

	authorizedKey, err = loadAuthorizedKey("client.pub")
	if err != nil {
		log.Fatalf("Failed to load authorized public key: %v", err)
	}

	server := ssh.Server{
		Addr:             ":2222",
		HostSigners:      []ssh.Signer{privateKey},
		Handler:          handleConnection,
		PublicKeyHandler: publicKeyAuthHandler,
		// PasswordHandler:  passwordHandler,
	}

	log.Println("Starting SSH server on port 2222...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start SSH server:", err)
	}
}
