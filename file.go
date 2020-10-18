package ini

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"
)

// Listener interface used in watch function
type Listener interface{ listen(inifile string) }

// ListenFunc a default type to implement Listener interface
type ListenFunc func(string)

// Watch inifile if inifile is changed return new Config
func Watch(filename string, listener Listener) (*Config, error) {
	waitGrp := sync.WaitGroup{}
	waitGrp.Add(1)
	go func() {
		listener.listen(filename)
		waitGrp.Done()
	}()
	waitGrp.Wait()
	return Load(filename)
}

func (l ListenFunc) listen(filename string) {
	l(filename)
}

// Config represents a combination of one or more INI files in memory.
type Config struct {
	// To keep data in order.
	SectionList []string
	// Actual data is stored here.
	Sections map[string]*Section
}

// Load ini file, return error if exists
func Load(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	cfg, err := parse(buf)
	return cfg, err
}

func (cfg *Config) init() {
	cfg.Sections = make(map[string]*Section)
}

// NewSection creates a new section.
func (cfg *Config) NewSection(name string) error {
	if len(name) == 0 {
		return errors.New("empty section name")
	}
	if cfg.Sections[name] != nil {
		return fmt.Errorf("section(%s) name already exists", name)
	}
	cfg.SectionList = append(cfg.SectionList, name)
	cfg.Sections[name] = &Section{
		KeyVal: make(map[string]string),
	}
	return nil
}
