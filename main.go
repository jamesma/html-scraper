package main

import (
	"flag"
	"fmt"

	"github.com/jamesma/html-scraper/scrap"
)

var (
	target = flag.String("target", "soundcloud", "the target to scrap")
)

func main() {
	flag.Parse()

	switch *target {
	case "soundcloud":
		scrap.Soundcloud()
	// case "chamberorganizer":
	// 	scrap.ChamberOrganizer()
	default:
		fmt.Println("unsupported target:", *target)
	}
}
