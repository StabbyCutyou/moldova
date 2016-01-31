package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/StabbyCutyou/moldovan_slammer/moldova"
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
	didErr := false
	for i := 0; i < cfg.iterations; i++ {
		t, err := moldova.ParseTemplate(cfg.template)
		if err != nil {
			log.Print(err)
			didErr = true
		} else {
			os.Stdout.Write([]byte(t + "\n"))
		}
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
