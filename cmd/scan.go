/*
Copyright © 2021 terabyte3

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
	"strings"

	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
)

var stimeout int
var sports string

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		_, err := ping.NewPinger(args[0])
		if err != nil {
			return errors.New(color.RedString(err.Error()))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ports, _ := parsePorts(sports)
		for _, port := range ports {
			isopen, err := pingPort(args[0], port, stimeout)
			if isopen {
				fmt.Println(color.GreenString("[+] ") + "port " + color.CyanString(fmt.Sprint(port)) + " is open")
			} else if strings.Contains(err.Error(), "i/o timeout") {
				fmt.Println(color.CyanString("[i] ") + "port " + color.CyanString(fmt.Sprint(port)) + " couldn't be reached, could be firewall or congestion")
			} else if strings.Contains(err.Error(), "connection refused") {
				fmt.Println(color.RedString("[-] ") + "port " + color.CyanString(fmt.Sprint(port)) + " is closed")
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().IntVarP(&stimeout, "timeout", "t", 30, "timeout")
	scanCmd.Flags().StringVarP(&sports, "ports", "p", "22,23,554,3306,179,1080,445,5432,6379", "range of ports (e.g. 80,8080,443 | 10-22)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
