package main

import (
	"encoding/json"
	"log"
	"os"

	resource "github.com/jghiloni/helm-resource"
	"github.com/jghiloni/helm-resource/check"
)

func main() {
	var req check.Request
	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Fatal(err)
	}

	client := resource.NewClient(req.Source.SkipTLSValidation)

	resp, err := check.RunCommand(client, req)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		log.Fatal(err)
	}
}
