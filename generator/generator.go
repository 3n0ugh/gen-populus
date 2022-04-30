package generator

import (
	"encoding/csv"
	"fmt"
	"github.com/3n0ugh/gen-populus/internal/config"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

func Generate(cfg config.Config) error {
	femaleNameFile, closer, err := OpenFile(cfg.FemaleNameFile)
	defer closer()
	if err != nil {
		return err
	}
	femaleNames, err := csv.NewReader(femaleNameFile).ReadAll()

	maleNameFile, closer, err := OpenFile(cfg.MaleNameFile)
	defer closer()
	if err != nil {
		return err
	}
	maleNames, err := csv.NewReader(maleNameFile).ReadAll()

	lastnameFile, closer, err := OpenFile(cfg.LastnameFile)
	defer closer()
	if err != nil {
		return err
	}
	lastnames, err := csv.NewReader(lastnameFile).ReadAll()

	femaleCount := uint64(float64(cfg.TotalPopulation) * 0.5)
	maleCount := uint64(float64(cfg.TotalPopulation) * 0.5)

	var nameChan = make(chan string, cfg.TotalPopulation)

	var wg sync.WaitGroup
	wg.Add(2)
	// Generate female names
	go func(ch chan<- string) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := uint64(0); i < femaleCount; i++ {
			name := femaleNames[r.Intn(len(femaleNames))][0]
			lastname := lastnames[r.Intn(len(lastnames))][0]
			gender := "female"
			age := r.Intn(110) + 1
			ch <- fmt.Sprintf("%s\t%s\t%s\t%d", name, lastname, gender, age)
		}
		wg.Done()
	}(nameChan)

	// Generate male names
	go func(ch chan<- string) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := uint64(0); i < maleCount; i++ {
			name := maleNames[r.Intn(len(femaleNames))][0]
			lastname := lastnames[r.Intn(len(lastnames))][0]
			gender := "male"
			age := r.Intn(110) + 1
			ch <- fmt.Sprintf("%s\t%s\t%s\t%d", name, lastname, gender, age)
		}
		wg.Done()
	}(nameChan)

	wg.Wait()
	close(nameChan)

	for n := range nameChan {
		fmt.Println(n)
	}

	return nil
}

func OpenFile(filepath string) (*os.File, func(), error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open csv file")
	}

	closeFunc := func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close csv file: %v", err)
		}
	}

	return file, closeFunc, err
}
