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
	"fmt"
	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"time"
)
var depth int
var cinterval int
var ignoreRobots bool
// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:   "crawl <hostname>",
	Example: `crawl -d 4 google.com`,
	Short: "crawl a webpage",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(color.CyanString("[i] ")+ "beginning crawl...")
		c := colly.NewCollector(
			colly.AllowedDomains("www."+args[0], args[0]),
		)
		c.OnRequest(func(r *colly.Request) {
    		fmt.Println(color.GreenString("[+] ")+"found: "+r.URL.String())
		})
		c.OnHTML("a", func(e *colly.HTMLElement) {
			nextPage := e.Request.AbsoluteURL(e.Attr("href"))
			time.Sleep(time.Duration(cinterval)*time.Second)
   			c.Visit(nextPage)
		})
		c.IgnoreRobotsTxt = ignoreRobots
		c.Visit("https://" + args[0])
		fmt.Println("Done")
	},
}

func init() {
	crawlCmd.Flags().IntVarP(&depth, "depth", "d", 4, "recursion depth")
	crawlCmd.Flags().IntVarP(&cinterval, "interval", "i", 7, "wait time between page visits (in seconds)")
	crawlCmd.Flags().BoolVarP(&ignoreRobots, "robotstxt", "r", false, "ignore robots.txt (makes you an obvious attacker)")
	// crawlCmd.Flags().StringVarP(&depth, "interval", "i", 7, "wait time between page visits")
	rootCmd.AddCommand(crawlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// crawlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// crawlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
