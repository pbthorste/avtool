package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/pbthorste/avtool"
	"github.com/smallfish/simpleyaml"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/urfave/cli.v1"
)

var (
	version string
)

func main() {
	app := cli.NewApp()
	app.Name = "avtool"
	app.Version = version
	app.Usage = "Tool for working with Ansible Vault files"
	app.Commands = []cli.Command{
		{
			Name:      "view",
			Usage:     "<vaultfile.yml> [-p password] [-k (all|.|<keyname>)] [-f vault-password-file]",
			UsageText: "'<vaultfile.yml>' (view contents of given Ansible Vault file)",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "vault-password-file, f",
					Usage: "load password from 'VAULT_PASSWORD_FILE'",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "`password` to use",
				},
				cli.StringFlag{
					Name:  "key, k",
					Usage: "'key' to retrieve the value for; 'keys' to list just the key names",
					Value: "keys",
				},
			},
			Action: func(c *cli.Context) error {
				vaultPassword := c.String("vault-password-file")
				password := c.String("password")
				key := c.String("key")
				// 01. Input param validations
				err := validateCommandArgs(c)
				if err != nil {
					return err
				}
				vaultFileName, err := validateAndGetVaultFile(c)
				if err != nil {
					return err
				}
				// 02. password retrieval and decryption
				pw, err := retrievePassword(vaultPassword, password)
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				result, err := avtool.DecryptFile(vaultFileName, pw)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				// 03. print output
				if result != "" {
					secretsYamlBytes := []byte(result)
					secretsYaml, _ := simpleyaml.NewYaml(secretsYamlBytes)
					// parse file contents as yaml
					if secretsYaml != nil {
						// 01. if all names or values are given, retrieve all keys, values
						if key == "all" || key == "keys" {
							keyList, _ := secretsYaml.GetMapKeys()
							println(getDecoratedMessage(fmt.Sprintf("%d Key(s) in %s", len(keyList), vaultFileName)))
							for _, keyName := range keyList {
								println(getDecoratedMessage(keyName))
								if key == "all" {
									println(getYamlKeyValue(keyName, secretsYaml))
								}
							}
						} else {
							// 02. if specific key is given, retrieve only that key
							println(getDecoratedMessage(key))
							println(getYamlKeyValue(key, secretsYaml))
						}
					}
				} else {
					println(getDecoratedMessage(vaultFileName + " is empty!"))
				}
				return nil
			},
		},
		{
			Name:  "encrypt",
			Usage: "vaultfile.yml - encrypt contents of the given yaml file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "vault-password-file",
					Usage: "load password from `VAULT_PASSWORD_FILE`",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "`password` to use",
				},
			},
			Action: func(c *cli.Context) error {
				// 01. Input param validations
				err := validateCommandArgs(c)
				if err != nil {
					return err
				}
				vaultFileName, err := validateAndGetVaultFile(c)
				if err != nil {
					return err
				}
				//
				vaultPassword := c.String("vault-password-file")
				password := c.String("password")
				pw, err := retrievePassword(vaultPassword, password)
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				result, err := avtool.EncryptFile(vaultFileName, pw)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				err = ioutil.WriteFile(vaultFileName, []byte(result), 0600)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				getDecoratedMessage("Encryption successful")
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func validateCommandArgs(c *cli.Context) (err error) {
	if !c.Args().Present() {
		cli.ShowSubcommandHelp(c)
		return cli.NewExitError(errors.New(getDecoratedMessage(
			"ERROR: Empty or Invalid inputs! Please ref. to usage instructions!")), 2)
	}
	return nil
}

func validateAndGetVaultFile(c *cli.Context) (filename string, err error) {
	filename = strings.TrimSpace(c.Args().First())
	if filename == "" {
		cli.ShowSubcommandHelp(c)
		return filename, cli.NewExitError(errors.New(getDecoratedMessage(
			"ERROR: Filename not specified!  Please ref. to usage instructions!")), 2)
	} else {
		if fileInfo, err := os.Stat(filename); os.IsNotExist(err) {
			cli.ShowSubcommandHelp(c)
			return filename, cli.NewExitError(errors.New(getDecoratedMessage(
				"ERROR: file "+filename+" "+"doesn't exist!")), 2)
		} else {
			if fileInfo.IsDir() {
				cli.ShowSubcommandHelp(c)
				return filename, cli.NewExitError(errors.New(getDecoratedMessage(
					"ERROR: file "+filename+" is a "+"directory!")), 2)
			}
		}
	}
	// return filename, error on error; nil if no error;
	return filename, nil
}

func getYamlKeyValue(keyName string, secretsYaml *simpleyaml.Yaml) (secretValue string) {
	secretValue, keyErr := secretsYaml.Get(keyName).String()
	if keyErr != nil {
		secretValue = "Key " + keyName + " doesn't exist! " + keyErr.Error()
	}
	return secretValue
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

	return readPassword()
}

/*
Reads password from stdin without showing what was entered.
*/
func readPassword() (password string, err error) {
	println(getDecoratedMessage("Enter password: "))
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		err = errors.New("ERROR: could not input password, " + err.Error())
		return
	}
	password = string(bytePassword)
	return
}

func getDecoratedMessage(messageIn string) (messageOut string) {
	return fmt.Sprintf(strings.Repeat(".", 4) + " " + messageIn + " " +
		strings.Repeat(".", 80-len(messageIn)-6))
}
