package csv

import (
	"encoding/csv"
	"github.com/pkg/errors"
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

func (w *Writer) WriteAll(data [][]string) error {
	w.mutex.Lock()
	err := w.csvWriter.WriteAll(data)
	if err != nil {
		err = errors.Wrap(err, "failed to write CSV file")
	}
	w.mutex.Unlock()
	return err
}

func (w *Writer) Flush() {
	w.mutex.Lock()
	w.csvWriter.Flush()
	w.mutex.Unlock()
}
