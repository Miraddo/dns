package cli

import "flag"

var (
	Domain string
	Query  string
)

// init ...
func init() {
	flag.StringVar(&Domain, "domain", "", "add domain name")
	flag.StringVar(&Query, "query", "A", "set query")
	flag.Parse()
}
