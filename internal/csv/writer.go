package csv

import (
	"encoding/csv"
	"os"
	"sync"
)

type Writer struct {
	mutex     *sync.Mutex
	csvWriter *csv.Writer
}

func NewCSVWriter(csvFile *os.File) (*Writer, error) {
	w := csv.NewWriter(csvFile)
	return &Writer{csvWriter: w, mutex: &sync.Mutex{}}, nil
}

func (w *Writer) WriteAll(data [][]string) {
	w.mutex.Lock()
	w.csvWriter.WriteAll(data)
	w.mutex.Unlock()
}

func (w *Writer) Flush() {
	w.mutex.Lock()
	w.csvWriter.Flush()
	w.mutex.Unlock()
}
