package main

import (
	"fmt"
	"os"

	ds "gitlab.flowcloud.systems/creator-ops/go-deviceserver-client"
)

func main() {
	token, err := ds.TokenFromPSK(os.Getenv("DEVICESERVER_PSK"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println(token)
}
