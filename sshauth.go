package main

import (
	"log"
	"os"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

var authorizedKey gossh.PublicKey

func loadAuthorizedKey(filename string) (gossh.PublicKey, error) {
	publicBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	authorizedKey, _, _, _, err := gossh.ParseAuthorizedKey(publicBytes)
	if err != nil {
		return nil, err
	}

	return authorizedKey, nil
}

func passwordHandler(ctx ssh.Context, password string) bool {
	var isOK bool = ctx.User() == "yashi" && password == "123456"
	log.Println(ctx.RemoteAddr(), l("AUTH"), l("PWD"), isOK)
	return isOK
}

func publicKeyAuthHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	var isOK bool = keysEqual(key, authorizedKey)
	log.Println(ctx.RemoteAddr(), l("AUTH"), l("CERT"), isOK)
	return isOK
}

func keysEqual(a, b gossh.PublicKey) bool {
	return a.Type() == b.Type() && string(a.Marshal()) == string(b.Marshal())
}
