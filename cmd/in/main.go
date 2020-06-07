package main

import (
	"encoding/json"
	"log"
	"os"

	resource "github.com/jghiloni/helm-resource"
	"github.com/jghiloni/helm-resource/in"
)

func main() {
	var req in.Request
	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Fatal(err)
	}

	client := resource.NewClient(req.Source.SkipTLSValidation)

	resp, err := in.RunCommand(os.Args[1], client, req)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		log.Fatal(err)
	}
}
