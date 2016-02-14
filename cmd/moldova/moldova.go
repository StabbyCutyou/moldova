package main

import (
	"bytes"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/StabbyCutyou/moldova"
)

type config struct {
	iterations int
	template   string
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	cs, err := moldova.BuildCallstack(cfg.template)
	if err != nil {
		log.Fatal(err)
	}
	didErr := false
	result := &bytes.Buffer{}
	for i := 0; i < cfg.iterations; i++ {
		err := cs.Write(result)
		if err != nil {
			log.Print(err)
			didErr = true
		} else {
			os.Stdout.Write([]byte(result.String() + "\n"))
		}
		result.Reset()
	}

	if didErr {
		os.Exit(1)
	}
}

func getConfig() (*config, error) {
	n := flag.Int("n", 1, "The number of times to generate a line of output. Cannot be set lower than 1")
	t := flag.String("t", "", "The template to generate results from")
	flag.Parse()
	if *n <= 0 {
		*n = 1
	}
	if *t == "" {
		return nil, errors.New("You must provide a template using the -t option")
	}

	return &config{iterations: *n, template: *t}, nil
}
