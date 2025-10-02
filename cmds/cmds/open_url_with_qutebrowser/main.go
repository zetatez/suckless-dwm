package main

import (
	"cmds/plugins"
	"flag"
)

func main() {
	prtUrl := flag.String("url", "www.google.com", "input url")
	flag.Parse()
	plugins.OpenUrlWithQutebrowser(*prtUrl)()
}
