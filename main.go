package main

import "mock/crawl"

func main() {
	c := crawl.NewURLFinder("https://bureauzigzag.nl", 2)
	o := c.Find()
	for _, u := range o {
		println(u)
	}
}
