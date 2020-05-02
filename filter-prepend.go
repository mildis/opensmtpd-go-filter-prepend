package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"strings"
)

var dec *mime.WordDecoder
var prefix string
var encprefix string
var forceEncode bool

func init() {
	flag.StringVar(&prefix, "prefix", "[*EXT*]", "Prepend subject with <prefix> if not already present")
	flag.BoolVar(&forceEncode, "encode", false, "Encode prefix whether subject is encoded or not")
	flag.Parse()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	dec = new(mime.WordDecoder)
	encprefix = mime.QEncoding.Encode("utf-8", prefix)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "config|ready") {
			RegisterFilter()
			log.Println("filter-prepend registered with " + prefix)
			if forceEncode {
				log.Println("filter-prepend will always encode prefix to " + encprefix)
			}
		} else {
			dataSplit := strings.Split(line, "|")
			if len(dataSplit) >= 8 {
				if dataSplit[4] == "data-line" {
					DoDataLine(dataSplit)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func RegisterFilter() {
	fmt.Println("register|filter|smtp-in|data-line")
	fmt.Println("register|report|smtp-in|link-disconnect")
	fmt.Println("register|report|smtp-in|link-connect")
	fmt.Println("register|ready")
}

func DoDataLine(dataSplit []string) {
	if strings.HasPrefix(strings.ToUpper(dataSplit[7]), "SUBJECT: ") {
		fmt.Printf("filter-dataline|%s|%s|%s\n", dataSplit[5], dataSplit[6], ProcessSubject(dataSplit[7:]))
	} else {
		fmt.Printf("filter-dataline|%s|%s|%s\n", dataSplit[5], dataSplit[6], strings.Join(dataSplit[7:], "|"))
	}
}

func ProcessSubject(s []string) string {
	result := ""
	subject := ""
	isEncoded := false

	rawsub := strings.SplitAfterN(strings.Join(s, "|"), ": ", 2)[1]

	if len(rawsub) < 8 || !strings.HasPrefix(rawsub, "=?") || !strings.HasSuffix(rawsub, "?=") || strings.Count(rawsub, "?") != 4 {
		subject = rawsub
	} else {
		isEncoded = true
		subject, _ = dec.Decode(rawsub)
	}

	if !strings.Contains(subject, prefix) {
		if isEncoded || forceEncode {
			result = encprefix + " "
		} else {
			result = prefix + " "
		}
	}

	return "Subject: " + result + rawsub
}
