package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"runtime"
	"strings"
)

var dec *mime.WordDecoder
var prefix string

func init() {
	flag.StringVar(&prefix, "prefix", "[*EXT*]", "Prepend subject with <prefix> if not already present")
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(1)

	scanner := bufio.NewScanner(os.Stdin)
	dec = new(mime.WordDecoder)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "config|ready") {
			registerFilter()
			log.Println("filter add-ext registered with " + prefix)
		} else {
			dataSplit := strings.Split(line, "|")
			if len(dataSplit) >= 8 {
				if dataSplit[4] == "data-line" {
					doDataLine(dataSplit)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func registerFilter() {
	fmt.Println("register|filter|smtp-in|data-line")
	fmt.Println("register|report|smtp-in|link-disconnect")
	fmt.Println("register|report|smtp-in|link-connect")
	fmt.Println("register|ready")
}

func doDataLine(dataSplit []string) {
	if strings.HasPrefix(strings.ToUpper(dataSplit[7]), "SUBJECT: ") {
		fmt.Printf("filter-dataline|%s|%s|%s\n", dataSplit[6], dataSplit[5], processSubject(dataSplit[7:]))
	} else {
		fmt.Printf("filter-dataline|%s|%s|%s\n", dataSplit[6], dataSplit[5], strings.Join(dataSplit[7:], "|"))
	}
}

func processSubject(s []string) string {
	result := ""
	subject := ""

	rawsub := strings.SplitAfterN(strings.Join(s, "|"), ": ", 2)[1]

	if len(rawsub) < 8 || !strings.HasPrefix(rawsub, "=?") || !strings.HasSuffix(rawsub, "?=") || strings.Count(rawsub, "?") != 4 {
		subject = rawsub
	} else {
		subject, _ = dec.Decode(rawsub)
	}

	if !strings.Contains(subject, prefix) {
		result = prefix + " "
	}

	return "Subject: " + result + rawsub
}
