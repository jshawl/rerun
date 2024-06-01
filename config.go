package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config map[string][]string

func parseConfig() Config {
	var conf Config

	dat, err := os.ReadFile("./rerun.toml")
	if err != nil {
		log.Fatal(err)
	}

	_, err = toml.Decode(string(dat), &conf)

	if err != nil {
		log.Fatal(err)
	}

	return conf
}
