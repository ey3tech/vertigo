/*
Copyright © 2021 terabyte3 <terabyte@terabyteis.me>, T33R3x

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
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var verbose bool
var rootCmd = &cobra.Command{
	Use: "vertigo",
	Short: color.CyanString(`
    	          __  _           
 _   _____  _____/ /_(_)___ _____ 
| | / / _ \/ ___/ __/ / __ '/ __ \
| |/ /  __/ /  / /_/ / /_/ / /_/ /
|___/\___/_/   \__/_/\__, /\____/ 
		    /____/
	`) + "\n\nvertigo is is a CLI application for retrieving information about computers.", // http://patorjk.com/software/taag/#p=display&f=Slant&t=vertigo

}

func Execute() {
	rootCmd.DisableAutoGenTag = true
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbosity (-v)")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// helper functions
func pingPort(hostname string, port int, timeout int) (bool, error) {
	addr := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Millisecond)

	if err != nil {
		return false, err
	}

	defer conn.Close()
	return true, nil
}

func parsePorts(ports string) ([]int, error) {
	var portList []int

	if strings.Contains(ports, ",") {
		pports := strings.Split(ports, ",") // parsed ports

		for _, port := range pports {
			p, err := strconv.Atoi(port)

			if err != nil {
				return nil, err

			}
			portList = append(portList, p)
		}
	} else if strings.Contains(ports, "-") {
		pports := strings.Split(ports, "-") // parsed ports
		if len(pports) != 2 {
			return nil, errors.New("invalid port range")
		}
		plow, _ := strconv.Atoi(pports[0])
		phigh, _ := strconv.Atoi(pports[1])

		for i := plow; i <= phigh; i++ {
			portList = append(portList, i)
		}
	} else {
		p, _ := strconv.Atoi(ports)
		portList = append(portList, p)
	}
	return portList, nil
}
