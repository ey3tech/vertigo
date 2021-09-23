## vertigo crawl

crawl a webpage

```
vertigo crawl <hostname> [flags]
```

### Examples

```
vertigo crawl google.com
```

### Options

```
  -c, --crawlers int       amount of crawlers (unrecommended without proxies, defaults to 1) (default 1)
  -d, --depth int          recursion depth (default 4)
  -h, --help               help for crawl
  -i, --interval int       wait time between page visits (in seconds, defaults to 7) (default 1)
  -p, --proxylist string   list of proxies to use
  -r, --robotstxt          ignore robots.txt (makes you an obvious attacker, defaults to false)
  -s, --domains            makes the crawler prioritize finding all different domains, hugely boosts performance
  -t, --timeout int        page visit timeout (in seconds, defaults to 3) (default 3)
```

### Options inherited from parent commands

```
  -v, --verbose   enable verbosity (-v)
```

### SEE ALSO

* [vertigo](vertigo.md) - vertigo is is a CLI application for retrieving information about computers.

