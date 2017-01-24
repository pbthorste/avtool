package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"strings"
	"github.com/pbthorste/avtool"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"io/ioutil"
	"errors"
)

func main() {
	app := cli.NewApp()
	app.Name = "avtool"
	app.Version = "1.0.0"
	app.Usage = "Tool for working with Ansible Vault files `asdf`"
	app.Commands = []cli.Command{
		{
			Name:    "view",
			Usage:   "vaultfile.yml - view contents of given Ansible Vault file",
			Flags:   []cli.Flag{
				cli.StringFlag{
					Name: "vault-password-file",
					Usage: "load password from `VAULT_PASSWORD_FILE`",
				},
				cli.StringFlag{
					Name: "password, p",
					Usage: "`password` to use",
				},
			},
			Action:  func(c *cli.Context) error {
				filename := strings.TrimSpace(c.Args().First())
				if filename == "" {
					return cli.NewExitError("ERROR: missing file argument", 2)
				}
				vaultPassword := c.String("vault-password-file")
				password := c.String("password")
				pw, err := retrievePassword(vaultPassword, password)
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				result, err := avtool.Decrypt(filename, pw)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				fmt.Println(result)
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func retrievePassword(vaultPasswordFile, passwordFlag string) (string, error) {
	if vaultPasswordFile != "" {
		if _, err := os.Stat(vaultPasswordFile); os.IsNotExist(err) {
			return "", errors.New("ERROR: vault-password-file, could not find: " + vaultPasswordFile)
		}
		pw, err := ioutil.ReadFile(vaultPasswordFile)
		if err != nil {
			return "", errors.New("ERROR: vault-password-file, " + err.Error())
		}
		return strings.TrimSpace(string(pw)), nil
	}
	if passwordFlag != "" {
		return passwordFlag, nil
	}

	fmt.Print("Enter password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", errors.New("ERROR: could not input password, " + err.Error())
	}
	return string(bytePassword), nil
}