package generator

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func Generate(population uint64) error {
	femaleNames := [][]string{{"Zeynep"}, {"Melek"}, {"Ayse"}}
	maleNames := [][]string{{"Cem"}, {"Serhat"}, {"Ali"}}
	lastNames := [][]string{{"Yilmaz"}, {"Karabulut"}, {"Arslan"}}

	femaleCount := uint64(float64(population) * 0.5)
	maleCount := uint64(float64(population) * 0.5)

	var nameChan = make(chan string, population)

	var wg sync.WaitGroup
	wg.Add(2)
	// Generate female names
	go func(ch chan<- string) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := uint64(0); i < femaleCount; i++ {
			name := femaleNames[r.Intn(len(femaleNames))][0]
			lastname := lastNames[r.Intn(len(lastNames))][0]
			gender := "female"
			ch <- fmt.Sprintf("%s\t%s\t%s", name, lastname, gender)
		}
		wg.Done()
	}(nameChan)

	// Generate male names
	go func(ch chan<- string) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := uint64(0); i < maleCount; i++ {
			name := maleNames[r.Intn(len(femaleNames))][0]
			lastname := lastNames[r.Intn(len(lastNames))][0]
			gender := "male"
			ch <- fmt.Sprintf("%s\t%s\t%s", name, lastname, gender)
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
