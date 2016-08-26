package main

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/urfave/cli"
	ds "gitlab.flowcloud.systems/creator-ops/go-deviceserver-client"
	"gitlab.flowcloud.systems/creator-ops/go-deviceserver-client/hateoas"
)

const keysCategory = "Key management"

var createKey = cli.Command{
	Name:      "create-key",
	Aliases:   []string{"ck"},
	Category:  keysCategory,
	Usage:     "Create a new key/secret",
	ArgsUsage: "<name>",
	Flags:     []cli.Flag{},
	Action: func(c *cli.Context) error {
		keyName := c.Args().Get(0)
		d, err := ds.Create(hateoas.Create(&hateoas.Client{
			EntryURL: deviceserverURL,
		}))
		if err != nil {
			return err
		}
		defer d.Close()

		credentials, err := ReadCredentials()
		if err != nil {
			return err
		}

		err = d.Authenticate(credentials)
		if err != nil {
			return err
		}

		key, err := d.CreateAccessKey(keyName)
		if err != nil {
			return err
		}

		fmt.Printf("Name:   %s\nKey:    %s\nSecret: %s\nSelf:   %s\n",
			key.Name,
			key.Key,
			key.Secret,
			key.Links.Self())
		return nil
	},
}

var listKeys = cli.Command{
	Name:      "list-keys",
	Aliases:   []string{"lk"},
	Category:  keysCategory,
	Usage:     "Lists the known access keys",
	ArgsUsage: " ",
	Flags:     []cli.Flag{},
	Action: func(c *cli.Context) error {
		d, err := ds.Create(hateoas.Create(&hateoas.Client{
			EntryURL: deviceserverURL,
		}))
		if err != nil {
			return err
		}
		defer d.Close()

		credentials, err := ReadCredentials()
		if err != nil {
			return err
		}

		err = d.Authenticate(credentials)
		if err != nil {
			return err
		}

		var previous *ds.AccessKeys = nil
		count := 0
		for {
			keys, err := d.GetAccessKeys(previous)
			if err != nil {
				return err
			}
			if keys == nil {
				break
			}
			for _, key := range keys.Items {
				fmt.Printf("[%d] '%s' = %s\n  %s\n\n", count, key.Name, key.Key, key.Links.Self())
				count++
			}
			previous = keys
		}

		return nil
	},
}

var deleteKey = cli.Command{
	Name:      "delete-key",
	Aliases:   []string{"dk"},
	Category:  keysCategory,
	Usage:     "Delete the specified key",
	ArgsUsage: "<key self URL>",
	Flags:     []cli.Flag{},
	Action: func(c *cli.Context) error {
		self := c.Args().Get(0)

		u, err := url.Parse(self)
		if err != nil {
			return err
		}
		du, err := url.Parse(deviceserverURL)
		if err != nil {
			return err
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return errors.New("Invalid scheme for self link")
		}
		if u.Host != du.Host {
			return errors.New("self link is not for this deviceserver")
		}

		d, err := ds.Create(hateoas.Create(&hateoas.Client{
			EntryURL: deviceserverURL,
		}))
		if err != nil {
			return err
		}
		defer d.Close()

		credentials, err := ReadCredentials()
		if err != nil {
			return err
		}

		err = d.Authenticate(credentials)
		if err != nil {
			return err
		}

		err = d.Delete(self)
		if err != nil {
			return err
		}

		return nil
	},
}
