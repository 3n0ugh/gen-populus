package main

import (
	"github.com/3n0ugh/gen-populus/generator"
	"github.com/3n0ugh/gen-populus/internal/config"
	"log"
	"os"
)

func main() {
	// Open or if not exists create output file
	file, err := os.OpenFile("data.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("failed to close file: %v", err)
		}
	}()

	// If file parameters are not specified, the default paths are used.
	cfg, err := config.NewConfig(
		1e7,
		"",
		"",
		"",
		file)
	if err != nil {
		log.Fatal(err)
	}

	err = generator.Generate(cfg)
	if err != nil {
		log.Println(err)
	}
}
