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
	"strings"
)

func main() {
	modePtr := flag.String("mode", "check", "Can be check etc.")
	srcPtr := flag.String("source", "https://www.raspberrypi.org/magpi-issues/", "The URL to get the issues from")

	flag.Parse()

	fmt.Println("mode: ", *modePtr)
	fmt.Println("src: ", *srcPtr)

	switch *modePtr {
	case "check":
		check(srcPtr)
	case "debug":
		debug(srcPtr)
	default:
		printUsage()
		os.Exit(1)
	}

}

func printUsage() {

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

func processLine(line string) {
	re := regexp.MustCompile(`.*<tr.*href="([^"]*\.pdf)".*tr>.*`)
	allMatches := re.FindAllStringSubmatch(line, -1)
	if len(allMatches) > 0 {
		fmt.Println(allMatches[0])
	}
}

func debug(srcPtr *string) {
	fileBytes, err := ioutil.ReadFile("issues")
	if err != nil {
		log.Fatal(err)
	}
	allLines := strings.Split(string(fileBytes), "\n")
	for _, line := range allLines {
		processLine(line)
	}
}
