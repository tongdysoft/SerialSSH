package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gliderlabs/ssh"
	i18n "github.com/kagurazakayashi/libNyaruko_Go/nyai18n"
	"github.com/tarm/serial"
)

var (
	serialPort *serial.Port
	serialMu   sync.Mutex
)

func l(t string) string {
	return i18n.GetMultilingualText(t)
}

func handleConnection(sshSession ssh.Session) {
	serialMu.Lock()
	defer serialMu.Unlock()
	if serialPort == nil {
		fmt.Fprintln(sshSession, l("SERIALINITERR"))
		return
	}

	// sshSession -> (serialPort, os.Stdout)
	var link1 string = fmt.Sprintf("%s -> (%s, %s)", l("SSHSESSION"), l("SERIALPORT"), l("STDOUT"))
	log.Println(l("PIPECREATE"), ":", link1)
	go func() {
		multiWriter := io.MultiWriter(serialPort, os.Stdout)
		_, err := io.Copy(multiWriter, sshSession)
		if err != nil {
			log.Println(l("SSHTOSERIALERR"), ":", link1, err)
		}
	}()

	// serialPort -> (sshSession, os.Stdout)
	var link2 string = fmt.Sprintf("%s -> (%s, %s)", l("SERIALPORT"), l("SSHSESSION"), l("STDOUT"))
	log.Println(l("PIPECREATE"), ":", link2)
	multiWriter := io.MultiWriter(sshSession, os.Stdout)
	_, err := io.Copy(multiWriter, serialPort)
	if err != nil {
		log.Println(l("SERIALTOSSHERR"), ":", link2, err)
	}
}

func main() {
	i18n.AutoSetLanguage()
	i18n.Language = strings.Split(i18n.Language, "-")[0]
	i18n.LoadLanguageFile("language.ini", false)
	log.Println(l("TITLE"), "v1.0")

	go readSerialPort()

	privateKeyPath, err := filepath.Abs("server.pem")
	if err != nil {
		log.Fatalf("%s: %s: %v", l("CERTLOADERR"), "server.pem", err)
	}
	privateKey, err := loadOrGenerateSSHKey(privateKeyPath)
	if err != nil {
		log.Fatalf("%s: %s: %v", l("CERTLOADERR"), privateKeyPath, err)
	} else {
		log.Println(l("CERTLOADED"), privateKeyPath, privateKey.PublicKey().Type())
	}

	authorizedKeyPath, err := filepath.Abs("client.pub")
	if err != nil {
		log.Fatalf("%s: %s: %v", l("CERTLOADERR"), "client.pub", err)
	}
	authorizedKey, err = loadAuthorizedKey(authorizedKeyPath)
	if err != nil {
		log.Fatalf("%s: %s: %v", l("CERTLOADERR"), authorizedKeyPath, err)
	} else {
		log.Println(l("CERTLOADED"), authorizedKeyPath, authorizedKey.Type())
	}

	var server ssh.Server = ssh.Server{
		Addr:        "127.0.0.1:2222",
		HostSigners: []ssh.Signer{privateKey},
		Handler:     handleConnection,
		// PublicKeyHandler: publicKeyAuthHandler,
		// PasswordHandler:  passwordHandler,
	}

	log.Println(l("STARTSSH"), server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("%s: %v", l("STARTSSHERR"), err)
	}
}
