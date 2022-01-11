package main

import (
	"github.com/danangkonang/oauth2-golang/app"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	app.Run()
}
