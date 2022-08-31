package main

import (
	"flag"
	"log"

	"github.com/rm3l/devfile-lib-comment-preservation-issue/pkg/devfile"
)

func main() {
	inputPathFlag := flag.String("input", "./devfile.yaml", "path to the Devfile YAML file")
	newMetadataNameFlag := flag.String("name", "", "value to set for the 'metadata.name' field in the Devfile")
	flag.Parse()

	p := *inputPathFlag
	n := *newMetadataNameFlag
	log.Printf("Setting 'metadata.name' to %q for devfile at %q\n", n, p)
	err := devfile.UpdateName(p, n)
	if err != nil {
		panic(err)
	}
}
