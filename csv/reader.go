package csv

import (
	"encoding/csv"
	"github.com/pkg/errors"
	"os"
)

type Reader struct {
	csvReader *csv.Reader
}

func NewCSVReader(csvFile *os.File) (*Reader, error) {
	r := csv.NewReader(csvFile)
	return &Reader{csvReader: r}, nil
}

func (r *Reader) ReadeAll() ([][]string, error) {
	dt, err := r.csvReader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read CSV file")
	}

	if len(dt) == 0 {
		return nil, errors.New("CSV file is empty")
	}
	return dt, nil
}
