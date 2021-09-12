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
	"fmt"
	"os"
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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "", "v", false, "enable verbosity (-v)")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
