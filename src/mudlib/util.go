package mudlib

import (
	"fmt"
	"os"
)

func removeStringFromList(s string, list *[]string) error {
	for i, n := range *list {
		if n == s {
			*list = append((*list)[:i], (*list)[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%q not found in list\n", s)
}

func loadBytes(file string) ([]byte, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	length, err := f.Seek(0, 2)
	if err != nil {
		return nil, err
	}

	if length == 0 {
		return nil, fmt.Errorf("Empty file: %s", file)
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	b := make([]byte, length)

	if b == nil {
		return b, fmt.Errorf("Failed to create byte slice of length %d", length)
	}

	_, err = f.Read(b)
	return b, err
}
