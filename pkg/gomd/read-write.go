package gomd

import (
	"fmt"
	"os"
)

func Write(fileName, text string) error {
	err := os.WriteFile(fileName, []byte(text), 0666)
	if err != nil {
		return fmt.Errorf("Error writing file %s: %w", fileName, err)
	}
	return nil
}

func Read(fileName string) ([]byte, error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Filename %s does not exist: %w", fileName, err)
		}
		return nil, fmt.Errorf("Error loading file %s: %w", fileName, err)
	}
	return bytes, nil
}

// TODO: handle permissions, more errors (eg can't open due to perms), and tests
