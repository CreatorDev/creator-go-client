package main

import (
	"fmt"
	"os"

	"gitlab.flowcloud.systems/creator-ops/go-deviceserver-client"
)

func main() {
	config := deviceserver.Config{
		PSK: os.Getenv("DEVICESERVER_PSK"),
	}
	client, _ := deviceserver.Create(&config)
	token, _ := client.TokenPSK()
	fmt.Println(token)
}
