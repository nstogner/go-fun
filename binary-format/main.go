package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nstogner/go-fun/binary-format/mps7"
)

func main() {
	f, err := mps7.ReadFile(os.Stdin)
	if err != nil {
		log.Fatalf("reading file: %s", err)
	}

	if err := f.Validate(); err != nil {
		log.Fatalf("invalid file: %s", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(f); err != nil {
		log.Fatalf("writing parsed file: %s", err)
	}
}
