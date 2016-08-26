package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/urfave/cli"
	ds "gitlab.flowcloud.systems/creator-ops/go-deviceserver-client"
)

var (
	deviceserverURL string
	credentialsFile string
	keyName         string
)

var keyNameFlag = cli.StringFlag{
	Name:  "key-name, k",
	Usage: "Specifies the name of a new key you're trying to create",
}

func ReadCredentials() (*ds.AccessKey, error) {
	if credentialsFile[:2] == "~/" {
		credentialsFile = strings.Replace(credentialsFile, "~", os.Getenv("HOME"), 1)
	}
	credFile, err := os.Open(credentialsFile)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(credFile)
	if err != nil {
		return nil, err
	}
	var key ds.AccessKey
	err = json.Unmarshal(buf, &key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func main() {
	app := cli.NewApp()
	app.Name = "ds-cli"
	app.Usage = "CLI interface to creatordev.io deviceserver"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "deviceserver-url, u",
			EnvVar:      "DEVICESERVER_URL",
			Value:       "https://deviceserver.creatordev.io",
			Destination: &deviceserverURL,
		},
		cli.StringFlag{
			Name:        "credentials, c",
			EnvVar:      "CREDENTIALS_FILE",
			Value:       "~/.ds-cli",
			Destination: &credentialsFile,
		},
	}
	app.Commands = []cli.Command{
		admin,     // hidden
		createOrg, // hidden
		createKey,
		deleteKey,
		listKeys,
	}

	app.Run(os.Args)
}
