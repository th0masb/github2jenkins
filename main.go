package main

import (
	"flag"
	"fmt"
	"github.com/th0masb/github2jenkins/conf"
	"log"
	"net/http"
	"regexp"
)

const yamlRx string = `^.*[.]ya?ml$`
const configFlag string = "config"
const configFlagDescription string = "Provides a path to the yaml file for configuring the server"

type args struct {
	configPath string
}

func main() {
	args := parseArgs()
	yamlMatcher := regexp.MustCompile(yamlRx)
	if !yamlMatcher.MatchString(args.configPath) {
		log.Fatalf("%s is not a path to a yaml file\n", args.configPath)
	} else {
		log.Printf("Config path set as %s\n", args.configPath)
	}
	config, err := conf.LoadConfig(args.configPath)
	if err == nil {
		log.Printf("Loaded %+v\n", config)
		//    http.HandleFunc("/path", handler)
		//    log.Fatal(http.ListenAndServe(":8080", nil))
	} else {
		log.Fatalf("Failed to load config: %s\n", err)
	}
}

func parseArgs() args {
	configPath := flag.String(configFlag, "", configFlagDescription)
	flag.Parse()
	if *configPath == "" {
		log.Fatalf("Must provide config path: github2jenkins --config /path/to/config.yaml\n")
	}
	return args{configPath: *configPath}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
