package ini

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
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

func TestParse(t *testing.T) {
	want := Config{
		SectionList: []string{"DEFAULT", "paths", "server"},
		Sections: map[string]*Section{
			"DEFAULT": &Section{
				KeyVal: map[string]string{"app_mode": "development"},
			},
			"paths": &Section{
				KeyVal: map[string]string{"data": "/home/git/grafana"},
			},
			"server": &Section{
				KeyVal: map[string]string{
					"protocol":       "http",
					"http_port":      "9999",
					"enforce_domain": "true"},
			},
		},
	}
	f, err := os.Open(testDir + "my.ini")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	cfg, err := parse(buf)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if !reflect.DeepEqual(want.SectionList, cfg.SectionList) {
		t.Errorf("wantSecList: %v, gotSecList: %v", want.SectionList, cfg.SectionList)
	}
	if !reflect.DeepEqual(want.Sections, cfg.Sections) {
		for _, secName := range want.SectionList {
			if !reflect.DeepEqual(want.Sections[secName], cfg.Sections[secName]) {
				t.Errorf("wantSection: %v, gotSection: %v", want.Sections[secName], cfg.Sections[secName])
			}
		}
		t.Errorf("want Sections not equal to got Sections")
	}
}
