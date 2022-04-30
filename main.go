package main

import (
	"encoding/json"
	"fmt"
	"github.com/3n0ugh/gen-populus/internal/config"
	"github.com/3n0ugh/gen-populus/pkg/generator"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
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

	since := time.Since(start)

	// Prints the configuration information
	conf, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		log.Fatalf("failed to marshal config: %v", err)
	}

	fmt.Printf("%s\n", conf)
	fmt.Printf("Generated %d data\n", cfg.TotalPopulation)
	fmt.Println("Time elapsed:", since)

}
