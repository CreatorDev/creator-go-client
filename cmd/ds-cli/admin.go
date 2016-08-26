package main

import (
	"fmt"

	"github.com/urfave/cli"
	ds "gitlab.flowcloud.systems/creator-ops/go-deviceserver-client"
	"gitlab.flowcloud.systems/creator-ops/go-deviceserver-client/hateoas"
)

var (
	deviceserverPSK string
)

var pskFlag = cli.StringFlag{
	Name:        "psk",
	EnvVar:      "DEVICESERVER_PSK",
	Destination: &deviceserverPSK,
	Usage:       "(required)",
}

var admin = cli.Command{
	Name:      "admin",
	Hidden:    true,
	Usage:     "Uses PSK to generate a JWT access_token for admin purposes",
	ArgsUsage: " ",
	Flags: []cli.Flag{
		pskFlag,
	},
	Action: func(c *cli.Context) error {
		token, err := ds.TokenFromPSK(deviceserverPSK)
		if err != nil {
			return err
		}
		fmt.Println(token)
		return nil
	},
}

var createOrg = cli.Command{
	Name:      "create-org",
	Hidden:    true,
	Usage:     "Uses PSK to create a new organisation key/secret",
	ArgsUsage: "<name>",
	Flags: []cli.Flag{
		pskFlag,
	},
	Action: func(c *cli.Context) error {
		keyName := c.Args().Get(0)
		d, err := ds.Create(hateoas.Create(&hateoas.Client{
			EntryURL: deviceserverURL,
		}))
		if err != nil {
			return err
		}
		defer d.Close()

		token, _ := ds.TokenFromPSK(deviceserverPSK)
		d.SetBearerToken(token)

		key, err := d.CreateAccessKey(keyName)
		if err != nil {
			return err
		}

		err = WriteCredentials(key)
		return err
	},
}
