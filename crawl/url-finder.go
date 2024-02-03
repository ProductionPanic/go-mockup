package crawl

import (
	"fmt"
	"mock/url"
	"strconv"
	"sync"
)
import playwright "github.com/playwright-community/playwright-go"

type URLFinder struct {
	startUrl url.URL
	found    []string
	maxDepth int
	browser  *playwright.Browser
	pl       *playwright.Playwright
	db       *CrawlDB
}

func NewURLFinder(startUrl string, maxDepth int) *URLFinder {

	return &URLFinder{
		startUrl: *url.NewURL(startUrl),
		maxDepth: maxDepth,
		db:       DB(),
	}
}

func (f *URLFinder) Find() []string {

	pl, e := playwright.Run()
	if e != nil {
		panic(e)
	}
	browser, err := pl.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		panic(err)
	}
	f.browser = &browser
	f.pl = pl
	f.found = []string{}

	f.doUrl(f.startUrl, 0)

	(*f.browser).Close()
	f.pl.Stop()
	defer f.db.Close()
	return f.found
}

func (f *URLFinder) isFound(u url.URL) bool {
	for _, found := range f.found {
		if found == u.String() {
			return true
		}
	}
	return false
}

func (f *URLFinder) doUrl(u url.URL, depth int) {
	if f.isFound(u) {
		return
	}
	if depth > f.maxDepth {
		return
	}
	f.found = append(f.found, u.String())

	fmt.Println("found " + u.String())
	page, err := (*f.browser).NewPage()
	if err != nil {
		panic(err)
	}
	_, err = page.Goto(u.String())
	if err != nil {
		panic(err)
	}
	links, err := page.Locator("a").All()
	if err != nil {
		panic(err)
	}
	hrefs := []string{}
	for _, link := range links {
		href, err := link.GetAttribute("href")
		if err != nil {
			panic(err)
		}
		parsed := url.NewURL(href)
		if parsed == nil {
			continue
		}
		if parsed.SameHostAs(&f.startUrl) {
			if !f.isFound(*parsed) {
				err := f.db.InsertURL(parsed.String(), depth)
				if err != nil {
					return
				}
				err = f.db.InsertLink(u.String(), parsed.String())
				if err != nil {
					return
				}

				hrefs = append(hrefs, href)
			}
		}
	}

	e := page.Close()
	if e != nil {
		panic(e)
	}

	fmt.Println("found " + strconv.Itoa(len(hrefs)) + " links")
	var wg sync.WaitGroup
	batchCount := 3
	batchSize := len(hrefs) / batchCount
	batches := make([][]string, batchCount)
	for i := 0; i < batchCount; i++ {
		start := i * batchSize
		end := start + batchSize
		if i == batchCount-1 {
			end = len(hrefs)
		}
		batches[i] = hrefs[start:end]
	}
	for _, batch := range batches {
		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()
			for _, href := range batch {
				f.doUrl(*url.NewURL(href), depth+1)
			}
		}(batch)
	}
	wg.Wait()
}
