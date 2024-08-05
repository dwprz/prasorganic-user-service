package test

import (
	"log"
	"os"
	"path"
)

// ini untuk merubah working directory path saat menjalankan test supaya path nya berawal dari root

func init() {
	filename, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	dir := path.Join(filename, "../../../")

	err = os.Chdir(dir)
	if err != nil {
		log.Fatalf("failed to change current working directory: %v", err)
	}
}
