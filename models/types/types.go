package types

import (
	"os"
	"log"
	"errors"
	"database/sql/driver"
)

type PriceType float32

type DbFile struct {
	Path string
}

func (dbFile DbFile) AsFile() (*os.File, error) {
	if dbFile.Path == "" {
		return nil, nil
	}
	file, err := os.Open(dbFile.Path) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	return file, err
}

func (dbFile *DbFile) Scan(value any) error {

	str, ok := value.(string)
	if !ok {
		log.Fatal(ok)
		return errors.New("failed to unmarshal file value")
	}
	dbFile.Path = str
	return nil
}
func (dbFile DbFile) Value() (driver.Value, error) {
	_, err := dbFile.AsFile()
	if err != nil {
		return nil, errors.New("could not open file")
	}
	return dbFile.Path, nil
}