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
	_ "net"
	_ "time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// bruteCmd represents the brute command
var service string
var port int
var bruteCmd = &cobra.Command{
	Use:   "brute",
	Short: "crack a password",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(color.RedString("missing host to bruteforce"))
		}
		//conn, err := net.DialTimeout("tcp", args[0], time.Duration(timeout)*time.Millisecond)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("brute called")
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
