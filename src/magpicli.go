package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
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

		re := regexp.MustCompile(`.*<tr.*href="([^"]*\.pdf)".*tr>.*`)

		for {
			lineBytes, err := reader.ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println(re.FindAllString(string(lineBytes), -1))
			fmt.Println(re.FindAllStringSubmatch(string(lineBytes), -1))
		}

	}
}

func debug(srcPtr *string) {
	response, err := http.Get(*srcPtr)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		_, err := io.Copy(os.Stdout, response.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
}
