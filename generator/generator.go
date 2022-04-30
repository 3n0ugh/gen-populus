package generator

import (
	"encoding/csv"
	"fmt"
	"github.com/3n0ugh/gen-populus/internal/config"
	csvWriter "github.com/3n0ugh/gen-populus/internal/csv"
	"github.com/3n0ugh/snowflake"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

func Generate(cfg config.Config) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	femaleNameData, err := ReadCSV(cfg.FemaleNameFile)
	if err != nil {
		return err
	}

	maleNameData, err := ReadCSV(cfg.MaleNameFile)
	if err != nil {
		return err
	}

	lastnameData, err := ReadCSV(cfg.LastnameFile)
	if err != nil {
		return err
	}

	femaleCount := uint64(float64(cfg.TotalPopulation)*cfg.Gender.Female) + 1
	maleCount := uint64(float64(cfg.TotalPopulation) * cfg.Gender.Male)

	childCount := uint64(float64(cfg.TotalPopulation) * cfg.Elderly.Child)
	elderCount := uint64(float64(cfg.TotalPopulation) * cfg.Elderly.Elder)
	teenCount := cfg.TotalPopulation - childCount - elderCount

	node, err := snowflake.NewNode(1, 1)
	if err != nil {
		return errors.Wrap(err, "failed to create snowflake node")
	}

	type nameChanModel struct {
		Name, Gender string
	}

	type elderlyChanModel struct {
		Age       int
		Birthdate string
	}

	var nameChan = make(chan nameChanModel, cfg.TotalPopulation)
	var elderlyChan = make(chan elderlyChanModel, cfg.TotalPopulation)
	var wg sync.WaitGroup

	wg.Add(5)
	// Generate female names
	go func(nc chan<- nameChanModel) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := uint64(0); i < femaleCount; i++ {
			nc <- nameChanModel{Name: femaleNameData[r.Intn(len(femaleNameData))][0], Gender: "Female"}
		}
		wg.Done()
	}(nameChan)

	// Generate male names
	go func(nc chan<- nameChanModel) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := uint64(0); i < maleCount; i++ {
			nc <- nameChanModel{Name: maleNameData[r.Intn(len(maleNameData))][0], Gender: "Male"}
		}
		wg.Done()
	}(nameChan)

	// Generate child ages
	go func(ec chan<- elderlyChanModel) {
		defer wg.Done()
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := uint64(0); i < childCount; i++ {
			age := r.Intn(15) + 1 // Between 0-16
			year := time.Now().Year() - age
			month := r.Intn(11) + 1

			var day int
			switch month {
			case 1, 3, 5, 7, 8, 10, 12:
				day = r.Intn(30) + 1
			case 2:
				if year%4 == 0 {
					day = r.Intn(28) + 1
				} else {
					day = r.Intn(27) + 1
				}
			default:
				day = r.Intn(29) + 1
			}

			ec <- elderlyChanModel{
				Age:       age,
				Birthdate: fmt.Sprintf("%d.%d.%d", day, month, year),
			}
		}
	}(elderlyChan)

	// Generate teen ages
	go func(ec chan<- elderlyChanModel) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		defer wg.Done()
		for i := uint64(0); i < teenCount; i++ {
			age := r.Intn(40) + 16 // Between 16 and 56
			year := time.Now().Year() - age
			month := r.Intn(11) + 1

			var day int
			switch month {
			case 1, 3, 5, 7, 8, 10, 12:
				day = r.Intn(30) + 1
			case 2:
				if year%4 == 0 {
					day = r.Intn(28) + 1
				} else {
					day = r.Intn(27) + 1
				}
			default:
				day = r.Intn(29) + 1
			}

			ec <- elderlyChanModel{
				Age:       age,
				Birthdate: fmt.Sprintf("%d.%d.%d", day, month, year),
			}
		}
	}(elderlyChan)

	// Generate elder ages
	go func(ec chan<- elderlyChanModel) {
		defer wg.Done()
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := uint64(0); i < elderCount; i++ {
			age := r.Intn(55) + 56 // Between 56 and 111
			year := time.Now().Year() - age
			month := r.Intn(11) + 1

			var day int
			switch month {
			case 1, 3, 5, 7, 8, 10, 12:
				day = r.Intn(30) + 1
			case 2:
				if year%4 == 0 {
					day = r.Intn(28) + 1
				} else {
					day = r.Intn(27) + 1
				}
			default:
				day = r.Intn(29) + 1
			}

			ec <- elderlyChanModel{
				Age:       age,
				Birthdate: fmt.Sprintf("%d.%d.%d", day, month, year),
			}
		}
	}(elderlyChan)

	var dt = make([][]string, 0)
	w, err := csvWriter.NewCSVWriter(cfg.OutputFile)
	if err != nil {
		return errors.Wrap(err, "Failed to create csv writer")
	}

	for i := uint64(0); i < cfg.TotalPopulation; i++ {
		id, errId := node.Generate()
		if errId != nil {
			log.Println("failed to generate snowflake id")
		}

		n := <-nameChan
		e := <-elderlyChan
		lastname := lastnameData[r.Intn(len(lastnameData))][0]
		email := fmt.Sprintf("%s.%s_%d@3n0ugh.com", n.Name, lastname, r.Intn(10000))

		d := []string{id.String(), n.Name, lastname, email, strconv.Itoa(e.Age), e.Birthdate, n.Gender}
		dt = append(dt, d)

		if len(dt)%1e5 == 0 {
			wg.Add(1)
			go func(dt [][]string) {
				w.WriteAll(dt)
				w.Flush()
				wg.Done()
			}(dt)
			dt = nil
		}
	}

	wg.Wait()
	close(nameChan)
	return nil
}

func ReadCSV(filepath string) ([][]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open csv file")
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close csv file: %v", err)
		}
	}()

	dt, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read csv file")
	}

	if len(dt) == 0 {
		return nil, errors.New("csv file is empty")
	}
	return dt, nil
}
