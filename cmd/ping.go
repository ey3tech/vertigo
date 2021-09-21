/*
Copyright © 2021 terabyte3 <terabyte@terabyteis.me>

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
	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var count int
var interval int
var timeout int

func pingFunc(cmd *cobra.Command, args []string) error {
	p, _ := ping.NewPinger(args[0])
	err := p.Resolve()
	p.Interval = time.Duration(interval) * time.Millisecond
	p.Timeout = time.Duration(timeout) * time.Millisecond
	if err != nil && strings.Contains(err.Error(), "connect: no route to host") {
		return errors.New(color.RedString("couldn't find a route to the host"))
	}
	if err != nil && strings.Contains(err.Error(), "no such host") {
		return errors.New(color.RedString("i couldn't find that host at all, make sure there aren't any typos"))
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("resolved host")
	p.Count = count
	p.OnRecv = func(pkt *ping.Packet) {
		color.Green("[✓] recieved " + fmt.Sprint(pkt.Nbytes) + " bytes from " + string(pkt.Addr) + ": icmp_seq=" + fmt.Sprint(pkt.Seq) + " rtt=" + fmt.Sprint(pkt.Rtt))
	}
	p.SetPrivileged(false)
	fmt.Println("pinging", args[0], "with", count, "packets")
	err = p.Run()

	if err != nil && strings.Contains(err.Error(), "socket: permission denied") {
		p.SetPrivileged(true)
		err = p.Run()
	}
	if err != nil && strings.Contains(err.Error(), "socket: operation not permitted") {
		return errors.New(color.RedString("i don't have permission to send pings, try running me as sudo or running:\n\nsudo sysctl -w net.ipv4.ping_group_range=\"0 2147483647\""))
	}
	if err != nil {
		return errors.New(color.RedString(err.Error()))
	}
	if p.Statistics().PacketsRecv == 0 {
		return errors.New(color.RedString("request timed out"))
	}

	fmt.Println("-----------------------------------------------------------")
	color.Cyan("average latency: " + fmt.Sprint(p.Statistics().AvgRtt.Milliseconds()) + "ms")
	color.Cyan("min latency: " + fmt.Sprint(p.Statistics().MinRtt.Milliseconds()) + "ms")
	color.Cyan("max latency: " + fmt.Sprint(p.Statistics().MaxRtt.Milliseconds()) + "ms")
	return nil
}

var pingCmd = &cobra.Command{
	Use:   "ping <hostname> [count]",
	Short: "Check if a host is up and accepting basic incoming traffic.",
	Long:  `One of the most important things in`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(color.RedString("requires a host to ping"))
		}
		_, err := ping.NewPinger(args[0])
		if err != nil && strings.Contains(err.Error(), "connect: no route to host") {
			return errors.New(color.RedString("couldn't connect to the host"))
		}
		return nil
	},
	RunE: pingFunc,
}

func init() {
	pingCmd.Flags().IntVarP(&count, "count", "c", 4, "number of packets to send")
	pingCmd.Flags().IntVarP(&interval, "interval", "i", 1000, "interval between packet send")
	pingCmd.Flags().IntVarP(&timeout, "timeout", "t", 5000, "packet timeout duration")
	pingCmd.SuggestionsMinimumDistance = 1
	rootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
