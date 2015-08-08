package scrap

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	chamberOrganizerURLPrefix = "http://www.chamberorganizer.com/Calendar/moreinfo.php?eventid="
	eventNameToMatch          = "Event Name:"
	eventDescriptionToMatch   = "Description:"
	eventDateToMatch          = "Event Date:"
	eventTimeToMatch          = "Event Time:"
	eventLocationToMatch      = "Location:"
	eventContactPersonToMatch = "Contact Person:"
)

var (
	eventStringsToMatch = []string{
		eventNameToMatch,
		eventDescriptionToMatch,
		eventDateToMatch,
		eventTimeToMatch,
		eventLocationToMatch,
		eventContactPersonToMatch,
	}

	eventMatcher = func(n *html.Node, textToMatch string) bool {
		if n.DataAtom == atom.Font && n.Parent != nil && n.Parent.Parent != nil {
			parentSibling := n.Parent.PrevSibling
			if parentSibling != nil && parentSibling.FirstChild != nil {
				return strings.Contains(scrape.Text(parentSibling.FirstChild), textToMatch)
			}
		}
		return false
	}

	eventNameMatcher = func(n *html.Node) bool {
		return eventMatcher(n, eventNameToMatch)
	}

	eventDescriptionMatcher = func(n *html.Node) bool {
		return eventMatcher(n, eventDescriptionToMatch)
	}

	eventDateMatcher = func(n *html.Node) bool {
		return eventMatcher(n, eventDateToMatch)
	}

	eventTimeMatcher = func(n *html.Node) bool {
		return eventMatcher(n, eventTimeToMatch)
	}

	eventLocationMatcher = func(n *html.Node) bool {
		return eventMatcher(n, eventLocationToMatch)
	}

	eventContactPersonMatcher = func(n *html.Node) bool {
		return eventMatcher(n, eventContactPersonToMatch)
	}
)

func ensureAtMostOneElement(nodes []*html.Node, textToMatch string) {
	if len(nodes) > 1 {
		panic("encountered more than one html.Node match for: " + textToMatch)
	}
}

func extractEventDetails(root *html.Node) []*html.Node {
	eventNames := scrape.FindAll(root, eventNameMatcher)
	eventDescriptions := scrape.FindAll(root, eventDescriptionMatcher)
	eventDates := scrape.FindAll(root, eventDateMatcher)
	eventTimes := scrape.FindAll(root, eventTimeMatcher)
	eventLocations := scrape.FindAll(root, eventLocationMatcher)
	eventContacts := scrape.FindAll(root, eventContactPersonMatcher)

	// return nil if mandatory attributes are not found
	if len(eventNames) == 0 ||
		len(eventDates) == 0 ||
		len(eventContacts) == 0 {
		return nil
	}

	ensureAtMostOneElement(eventNames, eventNameToMatch)
	ensureAtMostOneElement(eventDescriptions, eventDescriptionToMatch)
	ensureAtMostOneElement(eventDates, eventDateToMatch)
	ensureAtMostOneElement(eventTimes, eventTimeToMatch)
	ensureAtMostOneElement(eventLocations, eventLocationToMatch)
	ensureAtMostOneElement(eventContacts, eventContactPersonToMatch)

	return []*html.Node{
		eventNames[0],
		eventDescriptions[0],
		eventDates[0],
		eventTimes[0],
		eventLocations[0],
		eventContacts[0],
	}
}

func eventDetailsToStrArr(eventDetails []*html.Node) []string {
	return []string{
		scrape.Text(eventDetails[0]),
		scrape.Text(eventDetails[1]),
		scrape.Text(eventDetails[2]),
		scrape.Text(eventDetails[3]),
		scrape.Text(eventDetails[4]),
		scrape.Text(eventDetails[5]),
		strings.TrimPrefix(
			scrape.Attr(eventDetails[5].FirstChild, "href"),
			"mailto:"),
	}
}

func writeCSVRecords(records [][]string, fileName string) {
	csvFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	for i := 0; i < len(records); i++ {
		record := records[i]
		if record != nil {
			writer.Write(records[i])
		}
	}
	writer.Flush()
}

// Scrap a single event for the specified eventID.
func scrapEvent(eventID int) []string {
	fullURL := chamberOrganizerURLPrefix + strconv.Itoa(eventID)
	fmt.Print("scraping... ", fullURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		panic(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	eventDetails := extractEventDetails(root)
	if eventDetails == nil {
		fmt.Println(" - invalid")
		return nil
	}

	fmt.Println(" - done")
	return eventDetailsToStrArr(eventDetails)
}

// Scrap events for the specified lowerEventID (inclusive) to
// upperEventID (exclusive).
func scrapEvents(lowerEventID int, upperEventID int) [][]string {
	delta := upperEventID - lowerEventID
	if delta < 0 {
		panic("lowerEventID (" + string(lowerEventID) +
			") is less than upperEventID (" + string(upperEventID) + ")")
	}

	ret := make([][]string, delta)
	for i := 0; i < delta; i++ {
		ret[i] = scrapEvent(lowerEventID + i)
	}
	return ret
}

// ChamberOrganizer scraps ChamberOrganizer. Scrap events for the specified
// lowerEventID (inclusive) to upperEventID (exclusive). Outputs to the
// specified file after scraping completes.
func ChamberOrganizer(lowerEventID int, upperEventID int, output string) {
	eventRecords := scrapEvents(lowerEventID, upperEventID)
	writeCSVRecords(eventRecords, output)
}
