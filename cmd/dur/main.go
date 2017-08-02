package main

import (
	"fmt"
	"os"
	"time"

	"github.com/raintank/dur"
)

func usage() {
	fmt.Println("dur <duration|time> <pattern>")
}

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(2)
	}
	if os.Args[2] == "" {
		usage()
		os.Exit(2)
	}
	patt := os.Args[2]

	mode := os.Args[1]
	switch mode {
	case "duration":
		val := dur.MustParseNDuration("input", patt)
		fmt.Println(val)
	case "time":
		now := time.Now()
		val := dur.MustParseDateTime(patt, time.Local, now, uint32(now.Unix()))
		fmt.Println(val)
	default:
		usage()
		os.Exit(2)

	}
}
