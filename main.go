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
	configFlag  string = "config"
	yamlRx      string = `^.*[.]ya?ml$`
	portFlag    string = "port"
	secretsFlag string = "secrets"
)

type args struct {
	configPath  string
	secretsPath string
	serverPort  string
}

func main() {
	args := parseArgs()
	config, err := g2j.LoadConfig(args.configPath, args.secretsPath)
	if err != nil {
		log.Fatalf("Failed to load config: %s\n", err)
	}
	log.Printf("Loaded config %+v\n", config)

	hookHandler := ingress.NewHookHandler(config)

	http.HandleFunc("/", hookHandler.Handle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", args.serverPort), nil))
}

func parseArgs() args {
	configPath := flag.String(configFlag, "", "")
	serverPort := flag.String(portFlag, "", "")
	secretsPath := flag.String(secretsFlag, "", "")
	flag.Parse()
	yamlMatcher := regexp.MustCompile(yamlRx)
	if !yamlMatcher.MatchString(*configPath) {
		log.Fatalf("Must provide .yaml config path: github2jenkins --config /path/to/config.yaml\n")
	}
	log.Printf("Config path set as %s\n", *configPath)
	return args{
		configPath:  *configPath,
		serverPort:  *serverPort,
		secretsPath: *secretsPath,
	}
}
