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
	"github.com/go-ping/ping"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/AlecAivazis/survey/v2"
)

// bruteCmd represents the brute command
var qs = []*survey.Question{
	{
        Name:     "hostname",
        Prompt:   &survey.Input{Message: "who are we cracking today?"},
        Validate: func (ans interface{}) error {
			// check if the hostname is invalid
			if str, ok := ans.(string) ; !ok {
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
            Options: []string{"http(s)", "ssh", "ftp"},
            Default: "ssh",
        },
    },
    {
        Name: "confirm",
        Prompt: &survey.Select{
			Message: "ready?",
			Options: []string{color.RedString("[x] no, i'm too baby"), color.GreenString("[+] let's goooo")},
			Default: 1,
		},
    },
}
var bruteCmd = &cobra.Command{
	Use:   "brute",
	Short: "crack a password",
	RunE: func(cmd *cobra.Command, args []string) error {
		answers := struct {
			Hostname string
			Service string
			Confirm int
		}{}
		if answers.Confirm == 0 { // index of the options list above
			fmt.Println(color.CyanString("[i] ") + "canceled")
			return nil
		}
		err := survey.Ask(qs, &answers)
		if err != nil {
			return errors.New(color.RedString(err.Error()))
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
