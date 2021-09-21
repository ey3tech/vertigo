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
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
)

// bruteCmd represents the brute command
type response struct {
	Hostname string
	Service  string
	Confirm  int
	Port     int
	Username string
	ListPath string
}

var answers = response{}
var port string
var qs = []*survey.Question{
	{
		Name:   "hostname",
		Prompt: &survey.Input{Message: "who are we cracking today?"},
		Validate: func(ans interface{}) error {
			// check if the hostname is invalid
			if str, ok := ans.(string); !ok {
				return errors.New(color.RedString("invalid hostname"))
			} else if _, ok := ping.NewPinger(str); ok != nil {
				return errors.New(color.RedString("invalid hostname"))
			}
			return nil
		},
		Transform: survey.ToLower,
	},
	{
		Name: "service",
		Prompt: &survey.Select{
			Message: "what service are we attacking?",
			Options: []string{"http", "ssh", "ftp"},
			Default: "ssh",
		},
	},
}

var bruteCmd = &cobra.Command{
	Use:   "brute",
	Short: "crack a password",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := survey.Ask(qs, &answers)
		if err != nil {
			return errors.New(color.RedString(err.Error()))
		}
		switch s := answers.Service; s {
		case "http":
			port = "80"
		case "ssh":
			port = "22"
		case "ftp":
			port = "21"
		}
		var qs2 = []*survey.Question{
			{
				Name: "port",
				Prompt: &survey.Input{
					Message: "what port should I use for the attack?",
					Default: port,
				},
				Validate: func(ans interface{}) error {
					_, err := strconv.Atoi(ans.(string))
					return err
				},
				Transform: func(ans interface{}) (newAns interface{}) {
					p, _ := strconv.Atoi(ans.(string))
					return p
				},
			},
			{
				Name: "username",
				Prompt: &survey.Input{
					Message: "what username are we cracking?",
				},
			},
			{
				Name: "listpath",
				Prompt: &survey.Input{
					Message: "where's the password list located?",
					Suggest: func(toComplete string) []string {
						var dir []fs.FileInfo
						var err error
						var path []string
						if toComplete == "" {
							dir, err = ioutil.ReadDir(".")
						} else {
							dir, err = ioutil.ReadDir(strings.Split(toComplete, "/")[0])
							dirstr := strings.Split(toComplete, "/")[0]
							path = strings.Split(toComplete, "/")
							fmt.Println(dirstr)
						}
						// fmt.Println(fmt.Sprint(len(dir)))
						if err != nil {
							return nil
						}
						var suggestions []string
						for _, item := range dir {
							if toComplete == "" {
								if item.IsDir() {
									suggestions = append(suggestions, item.Name() + "/")
								} else {
									suggestions = append(suggestions, item.Name())
								}
							} else if strings.Contains(item.Name(), path[1]) {
								if item.IsDir() {
									suggestions = append(suggestions, item.Name() + "/")
								} else {
									suggestions = append(suggestions, item.Name())
								}
							}
						}
						return suggestions
					},
				},
				Validate: func(ans interface{}) error {
					str := ans.(string)
					if !fs.ValidPath(str) {
						return errors.New(color.RedString("invalid filepath!"))
					}
					file, err := os.Stat(str)
					if file == nil || file.IsDir() {
						return errors.New(color.RedString("invalid filepath!"))
					}
					return err
				},
			},
			{
				Name: "confirm",
				Prompt: &survey.Select{
					Message: "ready?",
					Options: []string{color.RedString("[x] no i'm baby"), color.GreenString("[+] let's goooo")},
				},
			},
		}
		err = survey.Ask(qs2, &answers)
		if err != nil {
			return errors.New(color.RedString(err.Error()))
		}
		if answers.Confirm == 0 { // index of the options list above
			fmt.Println(color.CyanString("[i] ") + "canceled")
			return nil
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(bruteCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bruteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bruteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
