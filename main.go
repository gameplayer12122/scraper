package main

import (
	"flag"
	"fmt"
	"os"
	"scraper/scraper"
	"strconv"
	"strings"
	"sync"
)

var Retries = flag.Int("retries", 3, "Amount of times to retry whenever there's a error.")
var Threads = flag.Int("threads", 20, "Amount of files to download at the same time.")

var savePath = flag.String("dir", "", "The folder to save images in.")
var tags = normalizeTagsSlice(flag.Args())
var enabledSites = make(map[string]struct{})

func init() {
	flag.Usage = func() { // Custom usage
		writer := flag.CommandLine.Output()

		fmt.Fprintln(writer, "Scraper v3.0.0 by Luna")
		fmt.Fprintf(writer, "Usage: %s [options] [tag1 tag2 ... tagN]\n", os.Args[0])
		fmt.Fprintln(writer, "Tags are not downloaded separately - All tags are combined into one search query.")
		fmt.Fprintln(writer, "\nOptions:")

		flag.PrintDefaults()
	}

	availableSites := ""
	for _, mod := range scraper.Mods {
		availableSites += " "
		availableSites += mod.Name
	}

	enabledSitesString := flag.String("sites", "e621", "Sites to use for scraping. Separate multiple with the ',' character. For all sites, use 'all'. Currently available:"+availableSites)

	flag.Parse()

	if *enabledSitesString == "all" {
		for _, mod := range scraper.Mods {
			enabledSites[mod.Name] = struct{}{}
		}
	}

	for _, enabledSite := range strings.Split(*enabledSitesString, ",") {
		enabledSites[enabledSite] = struct{}{}
	}

	if *savePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	err := os.MkdirAll(*savePath, 0744)
	if err != nil {
		fmt.Printf("Could not create %s: %v\n", *savePath, err)
		os.Exit(1)
	}

}

func loop(mod scraper.Mod) {
	ModPrintf(mod.Name, "Starting scraper")

	page := 1
	tries := 0

	for {
		ModDebugf(mod.Name, "Scraping page %d", page)

		posts, err := mod.Func(tags, strconv.Itoa(page), client)

		if tries >= *Retries {
			ModPrintf(mod.Name, "%d retries exceeded. Failing with error: %v", *Retries, err)
			return
		} else if err != nil {
			ModPrintf(mod.Name, "Got error while scraping: %v", err)
			tries++
			continue
		}

		tries = 0

		if len(posts) == 0 {
			ModDebugf(mod.Name, "len(posts) == 0, assuming end")
			break
		}

		ModDebugf(mod.Name, "got %d posts (page %d, last id %d)", len(posts), page, posts[len(posts)-1].ID)

		for _, p := range posts {
			if SeenMD5(p.File.MD5) {
				ModDebugf(mod.Name, "MD5 %s already seen, skipping...", p.File.MD5)
				continue
			}

			AddMD5(p.File.MD5)

			Download(DownloadFile{
				Filename: mod.Name + "_" + strconv.FormatUint(p.ID, 10) + "." + p.File.Extension,
				URL:      p.File.URL,
			})
		}

		page++
	}

	ModPrintf(mod.Name, "Finished scraper")
}

func main() {
	wg := sync.WaitGroup{}

	for _, mod := range scraper.Mods {
		if _, ok := enabledSites[mod.Name]; !ok {
			continue
		}

		wg.Add(1)
		go func(mod scraper.Mod) {
			loop(mod)
			wg.Done()
		}(mod)
	}

	wg.Wait()

	fmt.Println("Waiting for downloads to finish")
	WaitDownloadFinish()
}
