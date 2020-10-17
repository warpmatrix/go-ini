package ini

import (
	"fmt"
	"reflect"
	"testing"
)

const testDir = "testdata/"

func TestLoadErr(t *testing.T) {
	cases := []struct {
		filename string
		wantErr  error
	}{
		{testDir + "my.ini", nil},
		{testDir + "nonexist.ini", fmt.Errorf("open testdata/nonexist.ini: no such file or directory")},
		{testDir + "empty.ini", nil},
	}
	for _, c := range cases {
		_, err := Load(c.filename)
		if err == nil && c.wantErr == nil {
			continue
		}
		if err == nil || c.wantErr == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("wantErr: %v, gotErr: %v", c.wantErr, err)
		}
	}
}

func TestLoadFile(t *testing.T) {
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
	cfg, err := Load(testDir + "my.ini")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if !reflect.DeepEqual(want.SectionList, cfg.SectionList) {
		t.Errorf("wantSecList: %v, gotSecList: %v", want.SectionList, cfg.SectionList)
	}
	for _, secName := range want.SectionList {
		if !reflect.DeepEqual(want.Sections[secName], cfg.Sections[secName]) {
			t.Errorf("wantSections: %v, gotSections: %v", want.Sections[secName], cfg.Sections[secName])
		}
	}
}
