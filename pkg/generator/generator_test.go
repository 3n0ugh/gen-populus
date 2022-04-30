package generator

import (
	"github.com/3n0ugh/gen-populus/internal/config"
	"github.com/pkg/errors"
	"os"
	"testing"
)

// BenchmarkGenerator is a benchmark that tests the Generate function to generate 10 million data.
func BenchmarkGenerator(b *testing.B) {
	// Create a temp output file
	tempOutFile, err := os.CreateTemp(".", "data.*.csv")
	if err != nil {
		b.Fatalf("Failed to create temp file: %s", err)
	}

	defer func() {
		if err := os.Remove(tempOutFile.Name()); err != nil {
			b.Fatalf("Failed to remove temp file: %s", err)
		}
	}()

	var cfg config.Config
	var errs = errors.New("")

	for errs != nil {
		cfg, errs = config.NewConfig(
			1e7,
			"../../internal/data/female_name.csv",
			"../../internal/data/male_name.csv",
			"../../internal/data/lastname.csv",
			tempOutFile)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := Generate(cfg)
		if err != nil {
			b.Fatalf("failed to generate data: %v", err)
		}
	}
}
