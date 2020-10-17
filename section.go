package ini

import "fmt"

// Section in ini file with key-value
type Section struct {
	KeyVal map[string]string
}

func (sec *Section) newKeyVal(keyName string, value string) error {
	if _, ok := sec.KeyVal[keyName]; ok {
		return fmt.Errorf("key(%v) already exists", keyName)
	}
	sec.KeyVal[keyName] = value
	return nil
}
