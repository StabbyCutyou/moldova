// Package csvutil is used by me to help quickly turn spreadsheets of data and print
// them as a specific go struct. Do not import or rely on this package for any
// reason, it will likely change in backwards-incompatible ways, and is not intended
// for others to use
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func PrintNamesFromCSV(path string) {
	fmt.Println("Reading from " + path)
	f, err := os.Open(path)
	if err != nil {
		// See? Don't use this. I panic in here, like a huge jerk. Not for use!
		panic(err)
	}
	r := csv.NewReader(f)

	header, err := r.Read()
	if err != nil {
		panic(err)
	}

	data, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, row := range data {
		fmt.Print("&Name{English,spellings{")
		for i, field := range row {
			if strings.Index(field, "-") > 0 {
				parts := strings.Split(field, "-")
				field = parts[0]
			}
			fmt.Printf("%s: \"%s\"", header[i], field)
			if i < len(header)-1 {
				fmt.Print(", ")
			} else {
				fmt.Print("}},")
			}
		}
		fmt.Print("\n")
	}
}

func main() {
	PrintNamesFromCSV(os.Args[1])
}
