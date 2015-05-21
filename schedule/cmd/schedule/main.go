package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/syamp/schedule"
	"gopkg.in/yaml.v2"
)

func main() {
	// options
	var path string
	var level string
	var location string
	var date string
	var name string

	flag.StringVar(&path, "path", "./examples",
		"path to directory containing schedule YAML input files")
	flag.StringVar(&level, "level", "primary",
		"level to lookup schedule for")
	flag.StringVar(&date, "date", "",
		"evaluate for date - format is "+schedule.TIME_FULLY_QUALIFIED)
	flag.StringVar(&location, "location", "UTC",
		"location to evaluate date option")
	flag.StringVar(&name, "name", "",
		"Name of schedule to lookup")
	flag.Parse()

	if name == "" {
		fmt.Println("name is required")
		flag.Usage()
		os.Exit(1)
	}

	when := time.Now().UTC()
	if date != "" {
		parsed, err := schedule.ParseTimeInLocation(schedule.TIME_FULLY_QUALIFIED, date, location)

		if err != nil {
			fmt.Println("unable to parse date")
			flag.Usage()
			os.Exit(1)
		}
		when = parsed
	}

	data, err := ioutil.ReadFile(path + "/" + name + ".yaml")
	m := schedule.Schedule{}
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		log.Fatalf("Error unmarshalling yaml: %v", err)
	}
	matches, err := m.Get(level, when)
	if err != nil {
		log.Fatalf("Unable to get any matches: %v", err)
	}
	fmt.Println(strings.Join(matches, ","))
}
