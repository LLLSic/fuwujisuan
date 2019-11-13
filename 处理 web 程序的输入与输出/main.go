package main

import (
	"cloudgo-io/service"
	"os"
	flag "github.com/spf13/pflag"
)

const (
	PORT string = "8080"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}

	pPort := flag.StringP("port", "p", PORT, "http listening")
	flag.Parse()
	if len(*pPort) != 0 {
		port = *pPort
	}

	Ser := service.NewServer()
	Ser.Run(":" + port)
}