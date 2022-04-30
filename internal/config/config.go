package config

import (
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"time"
)

// Config struct holds the ratios of the different types of blocks.
type Config struct {
	TotalPopulation uint64

	Gender struct {
		Male   float64
		Female float64
	}
	Elderly struct {
		Child float64
		Teen  float64
		Elder float64
	}

	// Path of csv files
	FemaleNameFile string
	MaleNameFile   string
	LastnameFile   string

	// Path of the output file
	OutputFile *os.File
}

// IsValid checks validation of the Config struct.
func (c Config) IsValid() (map[string]string, bool) {
	var errs = make(map[string]string, 0)

	if c.TotalPopulation <= 0 {
		errs["TotalPopulation"] = "must be greater than 0"
	}
	if c.Gender.Male+c.Gender.Female != 1 {
		errs["Gender"] = "Total ratio must be 1"
	}
	if c.Elderly.Elder+c.Elderly.Teen+c.Elderly.Child != 1 {
		errs["Elderly"] = "Total ratio must be 1"
	}

	if c.OutputFile == nil {
		errs["OutputFile"] = "must be specified"
	}

	if len(errs) > 0 {
		return errs, false
	}
	return nil, true
}

// NewConfig creates a new config.
func NewConfig(population uint64, femaleNameCSVFile, maleNameCSVFile, lastnameCSVFile string, outputFile *os.File) (Config, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	maleRatio := (49.38 + r.Float64()*(50.23-49.38)) / 100
	femaleRatio := 1.0 - maleRatio

	childRatio := (18 + r.Float64()*(24.0-18)) / 100
	oldRatio := (13 + r.Float64()*(21.0-13)) / 100
	teenRatio := 1.0 - childRatio - oldRatio

	if femaleNameCSVFile == "" {
		femaleNameCSVFile = "./internal/data/female_name.csv"
	}
	if maleNameCSVFile == "" {
		maleNameCSVFile = "./internal/data/male_name.csv"
	}
	if lastnameCSVFile == "" {
		lastnameCSVFile = "./internal/data/lastname.csv"
	}

	cfg := Config{
		TotalPopulation: population,
		Gender: struct {
			Male   float64
			Female float64
		}{
			Male:   maleRatio,
			Female: femaleRatio,
		},
		Elderly: struct {
			Child float64
			Teen  float64
			Elder float64
		}{
			Child: childRatio,
			Teen:  teenRatio,
			Elder: oldRatio,
		},
		FemaleNameFile: femaleNameCSVFile,
		MaleNameFile:   maleNameCSVFile,
		LastnameFile:   lastnameCSVFile,
		OutputFile:     outputFile,
	}

	if errs, ok := cfg.IsValid(); !ok {
		return Config{}, errors.New(fmt.Sprintf("%#v", errs))
	}
	return cfg, nil
}
