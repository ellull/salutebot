package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ellull/salutebot/pkg/saluter"
	"github.com/spf13/pflag"
)

var (
	url      = pflag.StringP("url", "u", "", "Webhook URL")
	filename = pflag.StringP("filename", "f", "goodmorning.csv", "Salutes file")
)

func init() {
	pflag.Parse()
}

func main() {
	if *url == "" {
		fmt.Println("url is mandatory")
		pflag.Usage()
		os.Exit(1)
	}

	s, err := saluter.NewFileSaluter(*filename)
	if err != nil {
		log.Fatalln(err)
	}
	if err := s.Salute(*url); err != nil {
		log.Fatalln(err)
	}
}
