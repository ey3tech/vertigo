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
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
)

var depth int
var cinterval int
var ctimeout int
var camount int
var ignoreRobots bool
var proxylist string

// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:     "crawl <hostname>",
	Example: `crawl -d 4 google.com`,
	Short:   "crawl a webpage",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(color.CyanString("[i] ") + "preparing...")

		var proxies []string

		if proxylist != "" {

			// make sure file exists
			info, err := os.Stat(proxylist)
			if os.IsNotExist(err) {
				return errors.New(color.RedString("proxy list file does not exist"))
			}
			if info.IsDir() {
				return errors.New(color.RedString(proxylist + " is a directory"))
			}

			// read file
			proxystr, err := os.ReadFile(proxylist)
			if proxystr == nil || string(proxystr) == "" {
				return errors.New(color.RedString("proxy list file is empty"))
			}
			if err != nil {
				return errors.New(color.RedString(err.Error()))
			}
			proxies = strings.Split(string(proxystr), "\n") // a list of proxies
		}
		c := colly.NewCollector(
			colly.AllowedDomains("www."+args[0], args[0]),
		)
		if camount > 1 {
			c.Limit(&colly.LimitRule{
				Parallelism: camount,
			})
			//c.Async = true
		}
		c.SetRequestTimeout(time.Duration(ctimeout) * time.Second)
		c.IgnoreRobotsTxt = ignoreRobots
		c.OnRequest(func(r *colly.Request) {
			if strings.HasSuffix(r.URL.String(), "#") || r.URL.String() == "#" {
				r.Abort()
			}
		})
		c.OnHTML("a", func(e *colly.HTMLElement) {
			if proxies != nil {
				rproxy := rand.Intn(len(proxies))
				c.SetProxy(proxies[rproxy])
				fmt.Println(color.GreenString("[+] ") + "found " + e.Request.URL.String() + " using proxy " + proxies[rproxy])
			} else {
				fmt.Println(color.GreenString("[+] ") + "found: " + e.Request.URL.String())
			}
			time.Sleep(time.Duration(cinterval) * time.Second)
			nextPage := e.Request.AbsoluteURL(e.Attr("href"))
			c.Visit(nextPage)
		})
		fmt.Println(color.CyanString("[i] ") + "beginning crawl...")
		c.AllowURLRevisit = false
		c.Visit("https://" + args[0] + "/")
		// c.Wait()
		return nil
	},
}

func init() {
	crawlCmd.Flags().IntVarP(&depth, "depth", "d", 4, "recursion depth")
	crawlCmd.Flags().IntVarP(&cinterval, "interval", "i", 7, "wait time between page visits (in seconds, defaults to 7)")
	crawlCmd.Flags().IntVarP(&timeout, "timeout", "t", 3, "page visit timeout (in seconds, defaults to 3)")
	crawlCmd.Flags().IntVarP(&camount, "crawlers", "c", 1, "amount of crawlers (unrecommended without proxies, defaults to 1)")
	crawlCmd.Flags().BoolVarP(&ignoreRobots, "robotstxt", "r", false, "ignore robots.txt (makes you an obvious attacker, defaults to false)")
	crawlCmd.Flags().StringVarP(&proxylist, "proxylist", "p", "", "list of proxies to use")
	rootCmd.AddCommand(crawlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// crawlCmd.PersistentFlags().String("foo", "", "A help for foo")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// crawlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
