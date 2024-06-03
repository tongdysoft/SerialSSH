package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/tarm/serial"
	gossh "golang.org/x/crypto/ssh" // 使用别名导入
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

func generateSSHKey(filename string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	privateKeyFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer privateKeyFile.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return fmt.Errorf("failed to write private key to file: %v", err)
	}

	return nil
}

func generateECDSAKey(filename string) error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate ECDSA private key: %v", err)
	}

	privateKeyFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer privateKeyFile.Close()

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal ECDSA private key: %v", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return fmt.Errorf("failed to write private key to file: %v", err)
	}

	return nil
}

func loadOrGenerateSSHKey(filename string) (gossh.Signer, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Println("Private key file not found, generating a new one...")
		if err := generateSSHKey(filename); err != nil {
			return nil, err
		}
	}

	privateBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %v", err)
	}

	private, err := gossh.ParsePrivateKey(privateBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return private, nil
}

func main() {
	go readSerialPort()

	privateKey, err := loadOrGenerateSSHKey("host.key.pem")
	if err != nil {
		log.Fatalf("Failed to load or generate private key: %v", err)
	}

	server := ssh.Server{
		Addr:        ":2222",
		HostSigners: []ssh.Signer{privateKey},
		Handler:     handleConnection,
	}

	log.Println("Starting SSH server on port 2222...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start SSH server:", err)
	}
}
