package main

import (
	"fmt"

	"github.com/urfave/cli"
	ds "gitlab.flowcloud.systems/creator-ops/go-deviceserver-client"
	"gitlab.flowcloud.systems/creator-ops/go-deviceserver-client/hateoas"
)

var (
	deviceserverPSK string
	organisationID  int
)

var pskFlag = cli.StringFlag{
	Name:        "psk",
	EnvVar:      "DEVICESERVER_PSK",
	Destination: &deviceserverPSK,
	Usage:       "(required)",
}

var organisationFlag = cli.IntFlag{
	Name:        "org-id",
	Destination: &organisationID,
	Usage:       "If provided",
	Value:       0,
}

var adminToken = cli.Command{
	Name:      "admin-token",
	Hidden:    true,
	Usage:     "Uses PSK to generate a JWT access_token for admin purposes",
	ArgsUsage: " ",
	Flags: []cli.Flag{
		organisationFlag,
		pskFlag,
	},
	Action: func(c *cli.Context) error {
		token, err := ds.TokenFromPSK(deviceserverPSK, organisationID)
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
	Usage:     "Uses PSK to create a new organisation key/secret. If the specified organisation ID is zero, a new non-zero one will be automatically created.",
	ArgsUsage: "<name>",
	Flags: []cli.Flag{
		organisationFlag,
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

		token, _ := ds.TokenFromPSK(deviceserverPSK, organisationID)
		d.SetBearerToken(token)

		key, err := d.CreateAccessKey(keyName)
		if err != nil {
			return err
		}

		err = WriteCredentials(key)
		if err != nil {
			return err
		}
		fmt.Println("OK")
		return nil
	},
}
