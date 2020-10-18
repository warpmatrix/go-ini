package ini

import "fmt"

// Section in ini file with keyVal map
//
// Section 数据结构提供了 KeyVal 字段，可用通过 key 字符串映射到 value 字符串
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
