package ini

import (
	"fmt"
	"testing"
)

func TestParseKeyName(t *testing.T) {
	cases := []struct {
		line        string
		wantKeyName string
		wantOffset  int
		wantErr     error
	}{
		{"key1 = value", "key1", 6, nil},
		{" = value", "", 2, nil},
		{"key2  value", "", -1, fmt.Errorf("delimiter(%s) not found", keyValueDelim)},
	}
	for _, c := range cases {
		keyName, offset, err := parseKeyName(c.line)
		if c.wantErr == nil && err == nil {
			if c.wantKeyName != keyName {
				t.Errorf("wantKeyName: %v, gotKeyName: %v", c.wantKeyName, keyName)
			}
			if offset != c.wantOffset {
				t.Errorf("wantOffset: %v, gotOffset: %v", c.wantOffset, offset)
			}
		} else if err == nil || c.wantErr == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("wantErr: %v, gotErr: %v", c.wantErr, err)
		}
	}
}

func TestParseSecName(t *testing.T) {
	cases := []struct {
		line        string
		wantSecName string
		wantErr     error
	}{
		{"[secname]", "secname", nil},
		{"[]", "", fmt.Errorf("empty section name")},
		{"[uncloseSec", "", fmt.Errorf("unclosed section: %s", "[uncloseSec")},
	}
	for _, c := range cases {
		cfg := Config{}
		cfg.init()
		secName, err := parseSecName([]byte(c.line), &cfg)
		if c.wantErr == nil && err == nil {
			if c.wantSecName != secName {
				t.Errorf("wantSecName: %v, gotSecName: %v", c.wantSecName, secName)
			}
		} else if err == nil || c.wantErr == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("wantErr: %v, gotErr: %v", c.wantErr, err)
		}
	}
}

func TestParseValue(t *testing.T) {
	cases := []struct {
		line    string
		wantVal string
		wantErr error
	}{
		{"value", "value", nil},
		{"   value   ", "value", nil},
		{"  value " + CommentSym + "comment", "value", nil},
	}
	for _, c := range cases {
		value, err := parseValue(c.line)
		if c.wantErr == nil && err == nil {
			if c.wantVal != value {
				t.Errorf("wantVal: %v, gotVal: %v", c.wantVal, value)
			}
		} else if err == nil || c.wantErr == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("wantErr: %v, gotErr: %v", c.wantErr, err)
		}
	}
}
