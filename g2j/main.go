package main

import (
	"flag"
	"fmt"
	"github.com/th0masb/github2jenkins/g2j/conf"
	"log"
	"net/http"
	"regexp"
)

const yamlRx = `^.*[.]ya?ml$`
const configFlag = "config"
const configFlagDescription = "Provides a path to the yaml file for configuring the server"

func main() {
	configPath := flag.String(configFlag, "", configFlagDescription)
	flag.Parse()
	yamlMatcher := regexp.MustCompile(configFlag)
	if !yamlMatcher.MatchString(*configPath) {
		log.Fatalf("%s is not a path to a yaml file\n", *configPath)
	}
	log.Printf("Using configuration at %s\n", *configPath)
	config, err := conf.LoadConfig(configPath)
	if err == nil {
		log.Printf("Loaded %+v\n", config)
		//    http.HandleFunc("/path", handler)
		//    log.Fatal(http.ListenAndServe(":8080", nil))
	} else {
		log.Fatalf("Failed to load config: %s\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
