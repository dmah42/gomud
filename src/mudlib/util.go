package mudlib

import (
  "fmt"
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
