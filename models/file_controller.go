package models

import (
	"bufio"
	"fmt"
	"os"
)

var FILEPATH string = "token.txt"

func ReadFirstToken() (string, error) {
	file, err := os.Open(FILEPATH)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", nil
}


func DeleteFirstLine() error {
	file, err := os.Open(FILEPATH)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) == 0 {
		return fmt.Errorf("file is empty")
	}

	lines = lines[1:]

	file, err = os.Create(FILEPATH)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		fmt.Fprintln(file, line)
	}

	return nil
}
