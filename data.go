package main

import (
	"encoding/json"
	"os"
)

func saveDataFile(path string, exams TimedExams) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	err = json.NewEncoder(file).Encode(exams)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func readDataFile(path string) (exams TimedExams, err error) {
	// If dataFile does not exist, ignore error
	dataFile, err := os.ReadFile(path)

	if err != nil && !os.IsNotExist(err) {
		// If dataFile exists but there is an error, return it
		return nil, err
	} else if err != nil {
		// If dataFile does not exist, return empty array
		return nil, nil
	} else {
		err = json.Unmarshal(dataFile, &exams)
		if err != nil {
			return nil, err
		}
	}

	return exams, nil
}
