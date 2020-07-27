package main

import (
	"flag"
	"fmt"
	"github.com/th0masb/github2jenkins/g2j"
	"github.com/th0masb/github2jenkins/hook"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const configFlag string = "config"
const yamlRx string = `^.*[.]ya?ml$`
const portFlag string = "port"

type args struct {
	configPath string
	serverPort string
}

func main() {
	args := parseArgs()
	config, err := g2j.LoadConfig(args.configPath)
	if err == nil {
		log.Printf("Loaded %+v\n", config)
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", args.serverPort), nil))
	} else {
		log.Fatalf("Failed to load config: %s\n", err)
	}
}

func parseArgs() args {
	configPath := flag.String(configFlag, "", "")
	serverPort := flag.String(portFlag, "", "")
	flag.Parse()
	yamlMatcher := regexp.MustCompile(yamlRx)
	if !yamlMatcher.MatchString(*configPath) {
		log.Fatalf("Must provide .yaml config path: github2jenkins --config /path/to/config.yaml\n")
	} else {
		log.Printf("Config path set as %s\n", *configPath)
	}
	return args{configPath: *configPath, serverPort: *serverPort}
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Unable to read request body: %s\n", err)
		return
	}

	hook, err := hook.Parse(r.Header, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Unable to parse request body: %s %s\n", err, body)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Received hook: %+v\n", hook)
}
