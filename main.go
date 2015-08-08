package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jamesma/html-scraper/scrap"
)

var (
	targets = []string{
		"hackernews",
		"chamberorganizer",
	}

	listTargets = flag.Bool("list", false, "list all targets available")
	target      = flag.String("target", "hackernews", "the target to scrap")
	lowerID     = flag.Int("lowerID", 0, "the lower ID (inclusive)")
	upperID     = flag.Int("upperID", 0, "the upper ID (exclusive)")
	output      = flag.String("output", "output.csv", "the output file to write to")
)

func customUsageMessage() func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage: go main.go [flags]\n")
		flag.PrintDefaults()
	}
}

func listAllTargets() {
	fmt.Println("targets available:")
	for _, target := range targets {
		fmt.Println(target)
	}
}

func scrapeTarget(target *string) {
	switch *target {
	case "hackernews":
		scrap.HackerNews()
	case "chamberorganizer":
		scrap.ChamberOrganizer(*lowerID, *upperID, *output)
	default:
		fmt.Println("unsupported target:", *target)
	}
}

func main() {
	flag.Usage = customUsageMessage()
	flag.Parse()

	if *listTargets {
		listAllTargets()
		return
	}

	scrapeTarget(target)
}
