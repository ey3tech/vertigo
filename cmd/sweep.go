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
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
	"time"
)

func lookup(ipaddr string) {
	p, _ := ping.NewPinger(ipaddr)
	p.Count = 1
	p.Timeout = time.Duration(1 * time.Second)
	err := p.Run()
	if (err != nil) {
		if verbose {
			fmt.Println(color.RedString("[-] ") + ipaddr + " is unreachable")
		}
	} else {
		fmt.Println(color.CyanString("[+] ") +ipaddr+" is running")
	}
}

// sweepCmd represents the sweep command
var sweepCmd = &cobra.Command{
	Use:   "sweep <iprange>",
	Short: "ping all ip addresses in the given ip range",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New(color.RedString("requires exactly 1 argument"))
		}
		ip, ipnet, err := net.ParseCIDR(args[0])
		fmt.Println(ipnet.Network())
		if err != nil {
			return errors.New(color.RedString("invalid ip range"))
		}
		if ip.To4() == nil{
			return errors.New(color.RedString("invalid ip range"))
		}
		// p, err := ping.NewPinger(fmt.Sprint(ip))
		// go p.Run()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, ipnet, _ := net.ParseCIDR(args[0])
		mask := binary.BigEndian.Uint32(ipnet.Mask)
		start := binary.BigEndian.Uint32(ipnet.IP)

		// find the final address
		finish := (start & mask) | (mask ^ 0xffffffff)

		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, start)
		p, err := ping.NewPinger(fmt.Sprint(ip))
		p.Timeout = time.Duration(500)*time.Millisecond
		p.Count = 1
		p.OnRecv = func(pkt *ping.Packet) {
			fmt.Println(color.CyanString("[i]") + pkt.Addr, " is alive")
		}
		if err != nil {
			return errors.New(color.RedString(err.Error()))
		}
		for i := start; i <= finish; i++ {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, i)
			ipaddr, err := net.ResolveIPAddr("ip", fmt.Sprint(ip))
			if err != nil {
				return errors.New(color.RedString(err.Error()))
			}
			p.SetIPAddr(ipaddr)
			lookup(ipaddr.String())
			if err != nil {
				fmt.Println(color.RedString(err.Error()))
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sweepCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sweepCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sweepCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
