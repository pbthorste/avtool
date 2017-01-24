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
	"log"
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
				retrievePassword(vaultPassword, password)
				if password == "" {

					fmt.Print("Enter password: ")
					bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
					if err != nil {
						return cli.NewExitError("ERROR: could not input password, " + err.Error(), 2)
					}
					password = string(bytePassword)
				}
				fmt.Println("--" + password + "--")
				fmt.Println(avtool.Decrypt(filename, password))
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func retrievePassword(vaultPasswordFile, passwordFlag string) (string, error) {
	if vaultPasswordFile != "" {
		if _, err := os.Stat(vaultPasswordFile); os.IsNotExist(err) {
			return nil, error("ERROR: vault-password-file, could not find: " + vaultPasswordFile)
		}
		pw, err := ioutil.ReadFile(vaultPasswordFile)
		if err != nil {
			return nil, error("ERROR: vault-password-file, " + err.Error())
		}
		return strings.TrimSpace(pw), nil
	}
	if passwordFlag != "" {
		return passwordFlag, nil
	}

	fmt.Print("Enter password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, error("ERROR: could not input password, " + err.Error())
	}
	return string(bytePassword)
}