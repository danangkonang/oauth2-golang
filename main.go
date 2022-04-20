package main

import (
	"os"

	"github.com/danangkonang/oauth2-golang/app"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	os.Setenv("TZ", "Asia/Jakarta")
}

func main() {
	app.Run()
}
