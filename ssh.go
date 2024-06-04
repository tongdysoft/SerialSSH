package main

import "github.com/gliderlabs/ssh"

func passwordHandler(ctx ssh.Context, password string) bool {
	return ctx.User() == "yashi" && password == "123456"
}
