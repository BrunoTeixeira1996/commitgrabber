package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BrunoTeixeira1996/commitgrabber/internal"
)

func logic() error {
	var urlRepo = flag.String("url", "", ".repo.git link in HTTPS")
	var targetFileHash = flag.String("hash", "", "md5sum from the target file")
	var targetFileName = flag.String("fn", "", "name of the target file")
	flag.Parse()

	if *urlRepo == "" || *targetFileHash == "" || *targetFileName == "" {
		return fmt.Errorf("Please provide the url for the repository to look, the md5sum and the name of the target file")
	}

	if err := internal.GetCommit(*urlRepo, *targetFileHash, *targetFileName); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := logic(); err != nil {
		log.Fatal(err)
	}
}
