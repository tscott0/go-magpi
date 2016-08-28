package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Status int

const (
	HaveIssue     Status = iota // Already downloaded
	NewIssue                    // Ready to be downloaded - temperory
	WantedIssue                 // Matches user preferences
	UnwantedIssue               // Doesn't match user preferences
	FailedIssue                 // Failed
)

// A map of all the issues and their status
var allIssues = make(map[string]Status)

// Helper function to describe the status of an issue
func getStatusString(s Status) string {
	switch s {
	case HaveIssue:
		return "Have"
	case NewIssue:
		return "New"
	case WantedIssue:
		return "Wanted"
	case UnwantedIssue:
		return "Unwanted"
	case FailedIssue:
		return "Failed"
	default:
		return "Unknown"
	}
}

func main() {
	modePtr := flag.String("mode", "check", "check - download from the web\n    	debug - use a local issues file")
	srcPtr := flag.String("source", "https://www.raspberrypi.org/magpi-issues/", "The URL to get the issues from")
	bonusPtr := flag.Bool("bonus", false, "Download bonus issues")
	minIssue := flag.Int("minIssue", 0, "Only download if it's an newer issue than this")
	downloadDirPtr := flag.String("outputDir", "downloads", "The location to download issues to")

	flag.Parse()

	fmt.Println("mode:       ", *modePtr)
	fmt.Println("source:     ", *srcPtr)
	if *bonusPtr {
		fmt.Println("downloadBonus: true")
	} else {
		fmt.Println("downloadBonus: false")
	}
	fmt.Println("minIssue:   ", *minIssue)
	fmt.Println("outputDir:  ", *downloadDirPtr)
	fmt.Println()

	switch *modePtr {
	case "check":
		check(srcPtr)
	case "debug":
		debug(srcPtr)
		filterUnwanted(*bonusPtr, *minIssue)
		printIssueStatus(WantedIssue)
	default:
		flag.Usage()
		os.Exit(1)
	}

}

func printIssueStatus(s Status) {
	for issue, status := range allIssues {
		if status == s {
			fmt.Println(issue + " - " + getStatusString(status))
		}
	}
}

func check(srcPtr *string) {
	response, err := http.Get(*srcPtr)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()

		reader := bufio.NewReader(response.Body)

		for {
			lineBytes, err := reader.ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}
			processLine(string(lineBytes))
		}

	}
}

func debug(srcPtr *string) {
	fileBytes, err := ioutil.ReadFile("../issues")
	if err != nil {
		log.Fatal(err)
	}
	allLines := strings.Split(string(fileBytes), "\n")
	for _, line := range allLines {
		processLine(line)
	}
}

// Parse each line of the html and pull out issues.
// Issues are added to the issues map with a NewIssue Status
func processLine(line string) {
	// First check the line is a table row
	trRegex := regexp.MustCompile(`.*<tr.*href="([^"]*\.pdf)".*tr>.*`)
	allMatches := trRegex.FindAllString(line, -1)
	if len(allMatches) > 0 {
		// Pull out all .pdf entries. There should be two
		// one for the label, one for the link
		hrefRegex := regexp.MustCompile(`([^"><]*\.pdf)`)
		linkMatches := hrefRegex.FindAllString(line, -1)

		// Labels should always match but take the second element
		if len(linkMatches) == 2 {
			// Add if it doesn't exist
			currentIssue := linkMatches[1]
			if _, ok := allIssues[currentIssue]; !ok {
				allIssues[currentIssue] = 1
			}
		}
	}
}

// Iterate over all the issues and update the Status
func filterUnwanted(bonus bool, minIssue int) {
	for issue := range allIssues {
		// Change all issues to Wanted
		allIssues[issue] = WantedIssue

		if !isStandardIssue(issue) {
			// mark as unwanted if it's non-standard
			allIssues[issue] = UnwantedIssue
		} else if !isAboveMinimumIssue(issue, minIssue) {
			// or if it's above below the minimum issue number
			allIssues[issue] = UnwantedIssue
		}

	}
}

// Returns true if the issue has a standard name e.g. MagPi12.pdf
func isStandardIssue(issue string) bool {
	standardIssueRegex := regexp.MustCompile(`^MagPi[0-9]+.pdf$`)
	return standardIssueRegex.Match([]byte(issue))
}

func isAboveMinimumIssue(issue string, min int) bool {
	if isStandardIssue(issue) {
		issueNumRegex := regexp.MustCompile(`[0-9]+`)
		issueNumString := string(issueNumRegex.Find([]byte(issue)))
		issueNum, _ := strconv.Atoi(issueNumString)

		if issueNum < min {
			return false
		}
	}
	return true
}

func isAlreadDownloaded(issue string) bool {
	return false
	// TODO: Check for files that exist already so we don't download twice
}
