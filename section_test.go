package ini

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSecNewKeyVal(t *testing.T) {
	cases := []struct {
		key     string
		value   string
		wantErr error
	}{
		{"key1", "value", nil},
		{"key1", "value", fmt.Errorf("key(key1) already exists")},
		{"key2", "value", nil},
	}
	sec := Section{}
	sec.KeyVal = make(map[string]string)
	for _, c := range cases {
		err := sec.newKeyVal(c.key, c.value)
		if c.wantErr == nil && err == nil {
			continue
		} else if err == nil || c.wantErr == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("wantErr: %v, gotErr: %v", c.wantErr, err)
		}
	}
	want := Section{
		KeyVal: map[string]string{
			"key1": "value",
			"key2": "value",
		},
	}
	if !reflect.DeepEqual(sec.KeyVal, want.KeyVal) {
		t.Errorf("wantKeyVal: %v, gotKeyVal: %v", want.KeyVal, sec.KeyVal)
	}
}
