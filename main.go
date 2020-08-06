package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/th0masb/github2jenkins/g2j"
	"github.com/th0masb/github2jenkins/ingress"
)

const (
	configFlag string = "config"
	yamlRx     string = `^.*[.]ya?ml$`
	portFlag   string = "port"
)

type args struct {
	configPath string
	serverPort string
}

func main() {
	args := parseArgs()
	config, err := g2j.LoadConfig(args.configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %s\n", err)
	}
	log.Printf("Loaded config %+v\n", config)

	secrets, err := g2j.LoadSecrets(config.Secrets)
	if err != nil {
		log.Fatalf("Failed to load secrets: %s\n", err)
	}
	log.Printf("Loaded secrets.")

	hookHandler := ingress.NewHookHandler(config, secrets)

	http.HandleFunc("/", hookHandler.Handle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", args.serverPort), nil))
}

func parseArgs() args {
	configPath := flag.String(configFlag, "", "")
	serverPort := flag.String(portFlag, "", "")
	flag.Parse()
	yamlMatcher := regexp.MustCompile(yamlRx)
	if !yamlMatcher.MatchString(*configPath) {
		log.Fatalf("Must provide .yaml config path: github2jenkins --config /path/to/config.yaml\n")
	}
	log.Printf("Config path set as %s\n", *configPath)
	return args{configPath: *configPath, serverPort: *serverPort}
}
