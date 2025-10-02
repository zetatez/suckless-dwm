package main

import (
	"flag"

	"cmds/plugins"
)

func main() {
	prtUrl := flag.String("url", "www.google.com", "input url")
	flag.Parse()
	plugins.OpenUrlAsApp(*prtUrl)()
}
