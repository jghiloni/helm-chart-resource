package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"

	resource "github.com/jghiloni/helm-resource"
	"github.com/jghiloni/helm-resource/in"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.LUTC)
	var req in.Request

	var file *os.File
	var err error
	if file, err = ioutil.TempFile(os.TempDir(), "in-"); err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(io.TeeReader(os.Stdin, file))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Fatal(err)
	}
	file.Close()

	client := resource.NewClient(req.Source.SkipTLSValidation)

	resp, err := in.RunCommand(os.Args[1], client, req)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		log.Fatal(err)
	}
}
