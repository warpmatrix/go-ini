package ini

import (
	"fmt"
	"reflect"
	"testing"
)

const testDir = "testdata/"

func TestLoad(t *testing.T) {
	cases := []struct {
		filename string
		hasCfg   bool
		wantErr  error
	}{
		{testDir + "my.ini", true, nil},
		{testDir + "nonexist.ini", false, fmt.Errorf("open " + testDir + "nonexist.ini: no such file or directory")},
		{testDir + "empty.ini", true, nil},
	}
	for _, c := range cases {
		cfg, err := Load(c.filename)
		if err == nil && c.wantErr == nil {
			continue
		}
		if err == nil || c.wantErr == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("wantErr: %v, gotErr: %v", c.wantErr, err)
		}
		if (cfg != nil) != c.hasCfg {
			t.Errorf("want hasCfg: %v, got cfg: %v", c.hasCfg, cfg)
		}
	}
}

func TestConfigInit(t *testing.T) {
	cfg := Config{}
	cfg.init()
	if cfg.Sections == nil {
		t.Errorf("init cfg.sections error")
	}
}

func TestConfigNewSection(t *testing.T) {
	cases := []struct {
		name    string
		wantErr error
	}{
		{"sec1", nil},
		{"", fmt.Errorf("empty section name")},
		{"sec1", fmt.Errorf("section(sec1) name already exists")},
		{"sec2", nil},
	}
	cfg := Config{}
	cfg.init()
	for _, c := range cases {
		err := cfg.newSection(c.name)
		if c.wantErr == nil && err == nil {
			continue
		} else if err == nil || c.wantErr == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("wantErr: %v, gotErr: %v", c.wantErr, err)
		}
	}
	want := Config{
		SectionList: []string{"sec1", "sec2"},
		Sections: map[string]*Section{
			"sec1": &Section{KeyVal: map[string]string{}},
			"sec2": &Section{KeyVal: map[string]string{}},
		},
	}
	if !reflect.DeepEqual(cfg.SectionList, want.SectionList) {
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
